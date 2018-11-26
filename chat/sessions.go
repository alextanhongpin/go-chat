package chat

import (
	"crypto/rand"
	"encoding/base64"
	"sync"

	"github.com/gorilla/websocket"
)

// Session represents a user session. A user can have multiple sessions, since
// the user can be connected through multiple devices (web, mobile) and can
// open multiple tabs.
type Session struct {
	id     string
	conn   *websocket.Conn
	server string // The server the session resides on.
}

// NewSession returns a new websocket session.
func NewSession(conn *websocket.Conn) *Session {
	return &Session{
		// MD5 of timestamp + randomString(32) should give the right random string.
		id:   randomString(32),
		conn: conn,
		// ts:   time.Now(),
	}
}

func randomString(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(b)
}

// Conn returns the current session websocket connection.
func (s *Session) Conn() *websocket.Conn {
	return s.conn
}

// SessionID returns the id string of the session.
func (s *Session) SessionID() string {
	return s.id
}

// Sessions manages the session for the current socket.
type Sessions struct {
	mu       sync.RWMutex
	sessions map[string]*Session
}

// NewSessions returns a new session manager.
func NewSessions() *Sessions {
	return &Sessions{
		sessions: make(map[string]*Session),
	}
}

// Put inserts a new session.
func (s *Sessions) Put(sess *Session) {
	s.mu.Lock()
	s.sessions[sess.id] = sess
	s.mu.Unlock()
}

// Get returns a session by the given session id.
func (s *Sessions) Get(id string) *Session {
	s.mu.RLock()
	sess, _ := s.sessions[id]
	s.mu.RUnlock()
	return sess
}

// Delete removes a session from sessions.
func (s *Sessions) Delete(id string) {
	s.mu.Lock()
	delete(s.sessions, id)
	s.mu.Unlock()
}
