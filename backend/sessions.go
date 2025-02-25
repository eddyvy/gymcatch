package backend

import "sync"

// Sessions struct to store session data
type Sessions struct {
	mu       sync.RWMutex
	sessions map[string]string
}

// NewSessions creates a new Sessions instance
func NewSessions() *Sessions {
	return &Sessions{
		sessions: make(map[string]string),
	}
}

// Set adds or updates a session
func (s *Sessions) Set(sessionID, email string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[sessionID] = email
}

// Get retrieves a session by sessionID
func (s *Sessions) Get(sessionID string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	email, exists := s.sessions[sessionID]
	return email, exists
}

// Delete removes a session by sessionID
func (s *Sessions) Delete(sessionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, sessionID)
}
