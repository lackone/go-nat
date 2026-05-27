package main

import (
	"fmt"
	"go-nat/common"
	"io"
	"net"
	"time"

	"github.com/xtaci/smux"
)

type Client struct {
	ClientId   string // 客户端ID
	ServerAddr string // 服务器地址
}

func NewClient(clientId string, ServerAddr string) *Client {
	return &Client{
		ClientId:   clientId,
		ServerAddr: ServerAddr,
	}
}

func (client *Client) Run() {
	for {
		err := client.run()
		if err != nil && err != io.EOF {
			fmt.Printf("client run error: %s\n", err)
		}

		fmt.Printf("client reconnect\n")
		time.Sleep(time.Second)
	}
}

func (client *Client) run() error {
	conn, err := net.Dial("tcp", client.ServerAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	//发送handshake包
	req := &common.HandshakeReq{
		ClientId: client.ClientId,
	}
	buf, err := req.Encode()
	if err != nil {
		return err
	}

	conn.SetWriteDeadline(time.Now().Add(time.Second * 10))
	_, err = conn.Write(buf)
	conn.SetWriteDeadline(time.Time{})

	if err != nil {
		return err
	}

	//创建mux
	mux, err := smux.Client(conn, nil)
	if err != nil {
		return err
	}
	defer mux.Close()

	for {
		stream, err := mux.AcceptStream()
		if err != nil {
			return err
		}

		go client.handleStream(stream)
	}
}

// 处理stream
func (client *Client) handleStream(stream net.Conn) error {
	defer stream.Close()

	pp := &common.ProxyProtocol{}
	err := pp.Decode(stream)
	if err != nil {
		return err
	}

	//与本地建连接
	var localConn net.Conn
	switch pp.InternalProtocol {
	case "tcp":
		localConn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", pp.InternalIp, pp.InternalPort))
		if err != nil {
			return err
		}
		defer localConn.Close()

	default:
		return fmt.Errorf("invalid internal protocol: %s", pp.InternalProtocol)
	}

	//双向数据传输
	go func() {
		defer localConn.Close()
		defer stream.Close()

		io.Copy(localConn, stream)
	}()

	io.Copy(stream, localConn)

	return nil
}
