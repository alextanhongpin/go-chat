package chat

import (
	"sync"

	"github.com/gorilla/websocket"
)

// Session treats each socket connection as a new session and holds data.
type Session struct {
	mu   sync.RWMutex
	conn *websocket.Conn
	data map[string]interface{}
	// id is the unique session identifier.
	id SessionID
	// userID is the id of the user from the database.
	userID string
}

func NewSession(id string, ws *websocket.Conn) *Session {
	return &Session{
		conn: ws,
		data: make(map[string]interface{}),
		id:   SessionID(id),
	}
}

func (s *Session) Set(key string, value interface{}) {
	s.mu.Lock()
	s.data[key] = value
	s.mu.Unlock()
}

func (s *Session) Get(key string) interface{} {
	s.mu.RLock()
	data := s.data[key]
	s.mu.RUnlock()
	return data
}

func (s *Session) SessionID() string {
	s.mu.RLock()
	id := s.id
	s.mu.RUnlock()
	return id.String()
}

func (s *Session) Conn() *websocket.Conn {
	s.mu.RLock()
	conn := s.conn
	s.mu.RUnlock()
	return conn
}

// func (s *Session) Save(repo Repository) error {
//         return repo.Write(s.id, s)
// }
//
// func (s *Session) Remove(repo Repository) error {
//         repo.Delete(s.id)
// }
