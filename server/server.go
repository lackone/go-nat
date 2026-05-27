package main

import (
	"fmt"
	"go-nat/common"
	"net"
	"time"
)

type Server struct {
	ListenAddr string `json:"listen_addr"`
	sessionMgr *SessionManager
}

func NewServer(addr string, sessionMgr *SessionManager) *Server {
	s := &Server{
		ListenAddr: addr,
		sessionMgr: sessionMgr,
	}
	go s.checkOnline()
	return s
}

func (s *Server) ListenAndServe() error {
	// 监听TCP连接
	listen, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return err
	}
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			return err
		}

		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	req := &common.HandshakeReq{}
	err := req.Decode(conn)
	if err != nil {
		return
	}

	// 添加会话
	_, err = s.sessionMgr.AddSession(req.ClientId, conn)
	if err != nil {
		fmt.Errorf("AddSession failed: %v", err)
		return
	}

}

// checkOnline 检查会话是否在线
func (s *Server) checkOnline() {
	tick := time.NewTicker(time.Second * 3)
	defer tick.Stop()

	for range tick.C {
		s.sessionMgr.Range(func(k string, v *Session) bool {
			if v.Conn.IsClosed() {
				return false
			}
			return true
		})
	}
}
