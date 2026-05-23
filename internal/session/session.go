// Package session defines the core Session type used throughout sshtrace.
package session

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

// EventKind distinguishes input (keystrokes) from output (terminal data).
type EventKind string

const (
	EventInput  EventKind = "input"
	EventOutput EventKind = "output"
)

// Event represents a single captured terminal event within a session.
type Event struct {
	Timestamp time.Time `json:"timestamp"`
	Kind      EventKind `json:"kind"`
	Data      string    `json:"data"`
}

// Session holds metadata and events for a single SSH connection.
type Session struct {
	ID        string    `json:"id"`
	User      string    `json:"user"`
	ClientIP  string    `json:"client_ip"`
	StartedAt time.Time `json:"started_at"`
	EndedAt   time.Time `json:"ended_at,omitempty"`
	Events    []Event   `json:"events"`
	Tags      []string  `json:"tags,omitempty"`
}

// New creates a new Session with a random ID.
func New(user, clientIP string) (*Session, error) {
	id, err := randomID()
	if err != nil {
		return nil, err
	}
	return &Session{
		ID:        id,
		User:      user,
		ClientIP:  clientIP,
		StartedAt: time.Now(),
		Tags:      []string{},
	}, nil
}

// AddEvent appends an event to the session.
func (s *Session) AddEvent(kind EventKind, data string) {
	s.Events = append(s.Events, Event{
		Timestamp: time.Now(),
		Kind:      kind,
		Data:      data,
	})
}

// Close marks the session as finished.
func (s *Session) Close() { s.EndedAt = time.Now() }

// Duration returns how long the session lasted.
func (s *Session) Duration() time.Duration {
	if s.EndedAt.IsZero() {
		return time.Since(s.StartedAt)
	}
	return s.EndedAt.Sub(s.StartedAt)
}

func randomID() (string, error) {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
