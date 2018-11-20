package chat

import "sync"

type SessionID string

type Manager struct {
	mu       sync.RWMutex
	sessions map[SessionID]*Session
	repo     Repository
}

func NewManager(repo Repository) *Manager {
	return &Manager{
		sessions: make(map[SessionID]*Session),
		repo:     repo,
	}
}

func (s *Manager) Put(sess *Session) {
	s.mu.Lock()
	s.sessions[sess.id] = sess
	s.mu.Unlock()
}

func (s *Manager) Get(id string) *Session {
	s.mu.RLock()
	sess, _ := s.sessions[SessionID(id)]
	s.mu.RUnlock()
	return sess
}

func (s *Manager) Delete(id string) {
	s.mu.Lock()
	delete(s.sessions, SessionID(id))
	s.mu.Unlock()
}
