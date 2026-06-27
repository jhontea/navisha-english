package telegram

import (
	"sync"
	"time"
)

// PendingChallenge holds the active word challenge for a chat session
type PendingChallenge struct {
	ChallengeID        string
	IndonesianSentence string
	CorrectAnswer      string
	Topic              string
	CreatedAt          time.Time
}

// SessionStore manages per-chat pending challenges in memory
type SessionStore struct {
	mu       sync.RWMutex
	sessions map[int64]*PendingChallenge
}

func NewSessionStore() *SessionStore {
	s := &SessionStore{
		sessions: make(map[int64]*PendingChallenge),
	}
	// Background goroutine to clean up sessions older than 30 minutes
	go s.cleanup()
	return s
}

func (s *SessionStore) Set(chatID int64, challenge *PendingChallenge) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[chatID] = challenge
}

func (s *SessionStore) Get(chatID int64) (*PendingChallenge, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ch, ok := s.sessions[chatID]
	return ch, ok
}

func (s *SessionStore) Delete(chatID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, chatID)
}

func (s *SessionStore) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		s.mu.Lock()
		for chatID, ch := range s.sessions {
			if time.Since(ch.CreatedAt) > 30*time.Minute {
				delete(s.sessions, chatID)
			}
		}
		s.mu.Unlock()
	}
}
