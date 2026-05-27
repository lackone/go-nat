package main

import (
	"fmt"
	"net"
	"sync"

	"github.com/xtaci/smux"
)

type Session struct {
	ClientId string // 客户端ID
	//Conn     net.Conn // 连接
	//这里为什么不能用net.Conn,因为net.Conn一个连接只能处理一个数据流,而smux.Session可以打开多个数据流
	Conn *smux.Session
}

type SessionManager struct {
	sessionLock sync.RWMutex        // 会话锁
	Sessions    map[string]*Session // 会话映射
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		Sessions: make(map[string]*Session),
	}
}

// GetSessionByClientId 根据客户端ID获取会话
func (sm *SessionManager) GetSessionByClientId(clientId string) (net.Conn, error) {
	sm.sessionLock.RLock()
	defer sm.sessionLock.RUnlock()

	session, ok := sm.Sessions[clientId]
	if !ok {
		return nil, fmt.Errorf("session not found, clientId: %s", clientId)
	}

	//打开数据流
	stream, err := session.Conn.OpenStream()
	if err != nil {
		return nil, fmt.Errorf("session open stream failed, clientId: %s", clientId)

	}

	return stream, nil
}

// AddSession 添加会话
func (sm *SessionManager) AddSession(clientId string, conn net.Conn) (*Session, error) {
	sm.sessionLock.Lock()
	defer sm.sessionLock.Unlock()

	if _, ok := sm.Sessions[clientId]; ok {
		return nil, fmt.Errorf("session already exists, clientId: %s", clientId)
	}

	// 创建会话
	mux, err := smux.Server(conn, nil)
	if err != nil {
		return nil, err
	}

	sess := &Session{
		ClientId: clientId,
		Conn:     mux,
	}

	sm.Sessions[clientId] = sess

	return sess, nil
}

// CloseSession 关闭会话
func (sm *SessionManager) CloseSession(clientId string) {
	sm.sessionLock.Lock()
	defer sm.sessionLock.Unlock()

	sess, ok := sm.Sessions[clientId]
	if !ok {
		return
	}

	sess.Conn.Close()
	delete(sm.Sessions, clientId)
}

// Range 遍历所有会话
// 如果f返回false,则删除该会话
func (sm *SessionManager) Range(f func(k string, v *Session) bool) {
	sm.sessionLock.Lock()
	defer sm.sessionLock.Unlock()

	for k, v := range sm.Sessions {
		ok := f(k, v)
		if !ok {
			delete(sm.Sessions, k)
		}
	}
}
