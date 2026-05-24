// Package checkpoint provides session progress tracking so that
// long-running SSH sessions can be resumed or inspected mid-flight.
package checkpoint

import (
	"errors"
	"sync"
	"time"

	"github.com/sshtrace/sshtrace/internal/session"
)

// Checkpoint records the last known position within a session's event log.
type Checkpoint struct {
	SessionID  string    `json:"session_id"`
	EventIndex int       `json:"event_index"`
	SavedAt    time.Time `json:"saved_at"`
}

// Tracker maintains in-memory checkpoints keyed by session ID.
type Tracker struct {
	mu    sync.RWMutex
	store map[string]*Checkpoint
}

// New returns a new Tracker with an empty checkpoint store.
func New() *Tracker {
	return &Tracker{store: make(map[string]*Checkpoint)}
}

// Save records the current event index for the given session.
// Returns an error if the session is nil or the index is negative.
func (t *Tracker) Save(s *session.Session, index int) error {
	if s == nil {
		return errors.New("checkpoint: session must not be nil")
	}
	if index < 0 {
		return errors.New("checkpoint: event index must not be negative")
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.store[s.ID] = &Checkpoint{
		SessionID:  s.ID,
		EventIndex: index,
		SavedAt:    time.Now().UTC(),
	}
	return nil
}

// Load retrieves the checkpoint for the given session ID.
// Returns nil and false if no checkpoint exists.
func (t *Tracker) Load(sessionID string) (*Checkpoint, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	cp, ok := t.store[sessionID]
	return cp, ok
}

// Delete removes the checkpoint for the given session ID.
func (t *Tracker) Delete(sessionID string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.store, sessionID)
}

// All returns a snapshot of all stored checkpoints.
func (t *Tracker) All() []*Checkpoint {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make([]*Checkpoint, 0, len(t.store))
	for _, cp := range t.store {
		copy := *cp
		out = append(out, &copy)
	}
	return out
}
