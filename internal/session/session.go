package session

import (
	"time"

	"github.com/google/uuid"
)

// EventType represents the type of a session event.
type EventType string

const (
	EventTypeCommand EventType = "command"
	EventTypeOutput  EventType = "output"
	EventTypeConnect EventType = "connect"
	EventTypeDisconnect EventType = "disconnect"
)

// Event represents a single recorded event within an SSH session.
type Event struct {
	Timestamp time.Time `json:"timestamp"`
	Type      EventType `json:"type"`
	Data      string    `json:"data"`
}

// Session represents a recorded SSH session.
type Session struct {
	ID        string    `json:"id"`
	User      string    `json:"user"`
	RemoteIP  string    `json:"remote_ip"`
	StartedAt time.Time `json:"started_at"`
	EndedAt   *time.Time `json:"ended_at,omitempty"`
	Events    []Event   `json:"events"`
}

// New creates and returns a new Session with a unique ID.
func New(user, remoteIP string) *Session {
	return &Session{
		ID:        uuid.NewString(),
		User:      user,
		RemoteIP:  remoteIP,
		StartedAt: time.Now().UTC(),
		Events:    []Event{},
	}
}

// AddEvent appends a new event to the session's event log.
func (s *Session) AddEvent(eventType EventType, data string) {
	s.Events = append(s.Events, Event{
		Timestamp: time.Now().UTC(),
		Type:      eventType,
		Data:      data,
	})
}

// Close marks the session as ended.
func (s *Session) Close() {
	now := time.Now().UTC()
	s.EndedAt = &now
}

// Duration returns the duration of the session.
// If the session is still active, it returns the duration so far.
func (s *Session) Duration() time.Duration {
	if s.EndedAt != nil {
		return s.EndedAt.Sub(s.StartedAt)
	}
	return time.Since(s.StartedAt)
}
