package main

import (
	"fmt"
	"go-nat/common"
	"io"
	"net"
	"sync"
	"time"
)

var (
	WriteTimeout = time.Second * 10
)

type Listener struct {
	pp *common.ProxyProtocol

	tcpListener net.Listener    // TCP监听器
	sessionMgr  *SessionManager // 会话管理器

	closeCh   chan struct{} // 关闭通道
	closeOnce sync.Once     // 关闭一次
}

func NewListener(pp *common.ProxyProtocol, sessionMgr *SessionManager) *Listener {
	return &Listener{
		closeCh:    make(chan struct{}),
		sessionMgr: sessionMgr,
		pp:         pp,
	}
}

func (l *Listener) ListenAndServe() error {

	switch l.pp.PublicProtocol {
	case "tcp": // TCP协议
		return l.ListenAndServeTcp()
	default:
		return nil
	}

	return nil
}

func (l *Listener) ListenAndServeTcp() error {
	// 监听TCP连接
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", l.pp.PublicIp, l.pp.PublicPort))
	if err != nil {
		return err
	}
	defer listen.Close()

	l.tcpListener = listen // TCP监听器

	for {
		conn, err := listen.Accept()
		if err != nil {
			return err
		}

		go l.handleConn(conn)
	}
}

func (l *Listener) handleConn(conn net.Conn) {
	defer conn.Close()

	//查询session
	tunnelConn, err := l.sessionMgr.GetSessionByClientId(l.pp.ClientId)
	if err != nil {
		fmt.Errorf("session not found, clientId: %s", l.pp.ClientId)
		return
	}
	defer tunnelConn.Close()

	//封装proxy
	ppBody, err := l.pp.Encode()
	if err != nil {
		fmt.Errorf("encode proxy protocol failed, err: %v", err)
		return
	}

	//设置写超时
	tunnelConn.SetWriteDeadline(time.Now().Add(WriteTimeout))

	_, err = tunnelConn.Write(ppBody)

	//重置写超时
	tunnelConn.SetWriteDeadline(time.Time{})

	if err != nil {
		fmt.Errorf("write proxy protocol failed, err: %v", err)
		return
	}

	//双向数据传输
	go func() {
		defer tunnelConn.Close()
		defer conn.Close()

		io.Copy(tunnelConn, conn)
	}()

	io.Copy(conn, tunnelConn)
}

func (l *Listener) Close() {
	l.closeOnce.Do(func() {
		close(l.closeCh)

		if l.tcpListener != nil {
			l.tcpListener.Close() // 关闭TCP监听器
		}
	})
}
