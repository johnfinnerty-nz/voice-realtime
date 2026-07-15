package session

import (
	"sync"
	"time"
)

// State tracks a single client realtime session.
type State struct {
	ID        string
	Provider  string
	Model     string
	CreatedAt time.Time

	mu       sync.Mutex
	cancelled bool
}

// New creates session state.
func New(id, provider, model string) *State {
	return &State{
		ID:        id,
		Provider:  provider,
		Model:     model,
		CreatedAt: time.Now(),
	}
}

// Cancel marks the session as cancelled (barge-in).
func (s *State) Cancel() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cancelled = true
}

// IsCancelled reports whether response.cancel was received.
func (s *State) IsCancelled() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.cancelled
}

// ResetCancel clears the cancel flag after upstream acknowledges.
func (s *State) ResetCancel() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cancelled = false
}
