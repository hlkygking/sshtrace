// Package audit provides structured audit logging for SSH sessions,
// recording session lifecycle events to a persistent audit trail.
package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/sshtrace/sshtrace/internal/session"
)

// EventKind describes the type of audit event.
type EventKind string

const (
	EventSessionStarted EventKind = "session_started"
	EventSessionClosed  EventKind = "session_closed"
	EventCommandRun     EventKind = "command_run"
	EventAlertTriggered EventKind = "alert_triggered"
)

// Record is a single audit log entry.
type Record struct {
	Timestamp time.Time `json:"timestamp"`
	Kind      EventKind `json:"kind"`
	SessionID string    `json:"session_id"`
	User      string    `json:"user"`
	RemoteIP  string    `json:"remote_ip"`
	Detail    string    `json:"detail,omitempty"`
}

// Logger writes audit records to an io.Writer.
type Logger struct {
	w io.Writer
}

// New returns a Logger that writes JSON-encoded audit records to w.
func New(w io.Writer) *Logger {
	return &Logger{w: w}
}

// Log writes a single audit record for the given session and event kind.
func (l *Logger) Log(s *session.Session, kind EventKind, detail string) error {
	r := Record{
		Timestamp: time.Now().UTC(),
		Kind:      kind,
		SessionID: s.ID,
		User:      s.User,
		RemoteIP:  s.RemoteIP,
		Detail:    detail,
	}
	b, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("audit: marshal record: %w", err)
	}
	_, err = fmt.Fprintf(l.w, "%s\n", b)
	return err
}

// LogSession is a convenience wrapper that logs a session lifecycle event.
func (l *Logger) LogSession(s *session.Session, kind EventKind) error {
	return l.Log(s, kind, "")
}
