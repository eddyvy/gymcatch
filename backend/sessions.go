package backend

import (
	"sync"
)

// Sessions struct to store session data
type SessionsManager struct {
	mu       sync.RWMutex
	sessions map[string]*MegaCreds
}

var Sessions *SessionsManager

// NewSessions creates a new Sessions instance
func InitSessions() {
	Sessions = &SessionsManager{
		sessions: make(map[string]*MegaCreds),
	}
}

// Set adds or updates a session
func (s *SessionsManager) Set(sessionID string, creds *MegaCreds) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[sessionID] = creds
}

// Get retrieves a session by sessionID
func (s *SessionsManager) Get(sessionID string) (*MegaCreds, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	creds, exists := s.sessions[sessionID]
	return creds, exists
}

// Delete removes a session by sessionID
func (s *SessionsManager) Delete(sessionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, sessionID)
}
