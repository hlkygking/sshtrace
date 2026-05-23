// Package redact provides functionality to redact sensitive fields
// from session events before storage or export.
package redact

import (
	"regexp"
	"strings"

	"sshtrace/internal/session"
)

// DefaultFields contains common sensitive field names to redact.
var DefaultFields = []string{"password", "passwd", "secret", "token", "apikey", "api_key"}

// Redactor replaces sensitive values in session event data.
type Redactor struct {
	fields  []string
	pattern *regexp.Regexp
}

// New creates a Redactor. If fields is nil, DefaultFields are used.
func New(fields []string) (*Redactor, error) {
	if fields == nil {
		fields = DefaultFields
	}
	escaped := make([]string, len(fields))
	for i, f := range fields {
		escaped[i] = regexp.QuoteMeta(f)
	}
	patternStr := `(?i)(` + strings.Join(escaped, "|") + `)([=:\s]+)([^\s&;|]+)`
	re, err := regexp.Compile(patternStr)
	if err != nil {
		return nil, err
	}
	return &Redactor{fields: fields, pattern: re}, nil
}

// RedactEvent returns a copy of the event with sensitive data replaced.
func (r *Redactor) RedactEvent(e session.Event) session.Event {
	e.Data = r.pattern.ReplaceAllStringFunc(e.Data, func(match string) string {
		parts := r.pattern.FindStringSubmatch(match)
		if len(parts) < 4 {
			return match
		}
		return parts[1] + parts[2] + "[REDACTED]"
	})
	return e
}

// RedactSession returns a new session with all events redacted.
func (r *Redactor) RedactSession(s *session.Session) *session.Session {
	out := &session.Session{
		ID:        s.ID,
		User:      s.User,
		SourceIP:  s.SourceIP,
		StartTime: s.StartTime,
		EndTime:   s.EndTime,
	}
	for _, e := range s.Events {
		out.Events = append(out.Events, r.RedactEvent(e))
	}
	return out
}
