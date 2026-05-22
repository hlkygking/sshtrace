package session

import (
	"time"

	"github.com/google/uuid"
)

// Event represents a single recorded interaction within a session.
type Event struct {
	Timestamp time.Time
	Data      string
}

// Session holds metadata and events for a single SSH connection.
type Session struct {
	ID        string
	User      string
	RemoteIP  string
	StartedAt time.Time
	EndedAt   time.Time
	Events    []Event
}

// New creates a new Session for the given user and remote IP.
func New(user, remoteIP string) *Session {
	return &Session{
		ID:        uuid.New().String(),
		User:      user,
		RemoteIP:  remoteIP,
		StartedAt: time.Now(),
	}
}

// AddEvent appends an event to the session.
func (s *Session) AddEvent(e Event) {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now()
	}
	s.Events = append(s.Events, e)
}

// Close marks the session as ended.
func (s *Session) Close() {
	s.EndedAt = time.Now()
}

// Duration returns the elapsed time of the session.
// If the session is still open it measures from start until now.
func (s *Session) Duration() time.Duration {
	end := s.EndedAt
	if end.IsZero() {
		end = time.Now()
	}
	return end.Sub(s.StartedAt)
}
