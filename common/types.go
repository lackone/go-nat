package common

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
)

var (
	cmdProxyProtocol byte = 0x01
	cmdHandshakeReq  byte = 0x02
)

type ProxyProtocol struct {
	ClientId         string `json:"client_id"`         // 客户端ID
	PublicIp         string `json:"public_ip"`         // 公共IP
	PublicPort       int    `json:"public_port"`       // 公共端口
	PublicProtocol   string `json:"public_protocol"`   // 公共协议
	InternalIp       string `json:"internal_ip"`       // 内网IP
	InternalPort     int    `json:"internal_port"`     // 内网端口
	InternalProtocol string `json:"internal_protocol"` // 内网协议
}

//协议格式
//1byte version
//1byte cmd
//2byte length
//body

func (pp *ProxyProtocol) Encode() ([]byte, error) {
	buf := make([]byte, 4)
	buf[0] = 0
	buf[1] = cmdProxyProtocol

	body, err := json.Marshal(pp)
	if err != nil {
		return nil, err
	}

	binary.BigEndian.PutUint16(buf[2:4], uint16(len(body)))

	return append(buf, body...), nil
}

func (pp *ProxyProtocol) Decode(reader io.Reader) error {
	buf := make([]byte, 4)

	_, err := io.ReadFull(reader, buf)
	if err != nil {
		return err
	}

	if len(buf) != 4 {
		return fmt.Errorf("invalid handshake req, length: %d", len(buf))
	}

	cmd := buf[1]
	if cmd != cmdProxyProtocol {
		return fmt.Errorf("invalid handshake req, cmd: %s", cmd)
	}

	length := binary.BigEndian.Uint16(buf[2:4])

	body := make([]byte, length)
	_, err = io.ReadFull(reader, body)
	if err != nil {
		return err
	}
	if len(body) != int(length) {
		return fmt.Errorf("invalid handshake req, body length: %d", len(body))
	}

	err = json.Unmarshal(body, pp)
	if err != nil {
		return err
	}

	return nil
}

type HandshakeReq struct {
	ClientId string `json:"client_id"`
}

// Encode 编码握手请求
func (h *HandshakeReq) Encode() ([]byte, error) {
	buf := make([]byte, 4)
	buf[0] = 0
	buf[1] = cmdHandshakeReq

	body, err := json.Marshal(h)
	if err != nil {
		return nil, err
	}

	binary.BigEndian.PutUint16(buf[2:4], uint16(len(body)))

	return append(buf, body...), nil
}

// Decode 解码握手请求
func (h *HandshakeReq) Decode(reader io.Reader) error {
	buf := make([]byte, 4)

	_, err := io.ReadFull(reader, buf)
	if err != nil {
		return err
	}

	if len(buf) != 4 {
		return fmt.Errorf("invalid handshake req, length: %d", len(buf))
	}

	cmd := buf[1]
	if cmd != cmdHandshakeReq {
		return fmt.Errorf("invalid handshake req, cmd: %s", cmd)
	}

	length := binary.BigEndian.Uint16(buf[2:4])

	body := make([]byte, length)
	_, err = io.ReadFull(reader, body)
	if err != nil {
		return err
	}
	if len(body) != int(length) {
		return fmt.Errorf("invalid handshake req, body length: %d", len(body))
	}

	err = json.Unmarshal(body, h)
	if err != nil {
		return err
	}

	return nil
}
