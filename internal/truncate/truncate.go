// Package truncate provides utilities for limiting the size of session event
// data to prevent unbounded growth in captured output.
package truncate

import (
	"errors"
	"unicode/utf8"

	"github.com/sshtrace/sshtrace/internal/session"
)

const (
	// DefaultMaxBytes is the default maximum byte length for a single event's data.
	DefaultMaxBytes = 4096
	// DefaultMaxEvents is the default maximum number of events retained per session.
	DefaultMaxEvents = 1000
	// truncationMarker is appended when data is cut short.
	truncationMarker = "...[truncated]"
)

// Truncator limits event data size and event count within a session.
type Truncator struct {
	maxBytes  int
	maxEvents int
}

// New creates a Truncator with the given limits.
// maxBytes controls the maximum byte length of any single event's Data field.
// maxEvents controls the maximum number of events kept per session.
func New(maxBytes, maxEvents int) (*Truncator, error) {
	if maxBytes <= 0 {
		return nil, errors.New("truncate: maxBytes must be positive")
	}
	if maxEvents <= 0 {
		return nil, errors.New("truncate: maxEvents must be positive")
	}
	return &Truncator{maxBytes: maxBytes, maxEvents: maxEvents}, nil
}

// Apply enforces byte and event-count limits on the session in place.
// Events beyond maxEvents are dropped from the tail; event data exceeding
// maxBytes is trimmed and suffixed with a truncation marker.
func (t *Truncator) Apply(s *session.Session) {
	// Trim excess events first.
	if len(s.Events) > t.maxEvents {
		s.Events = s.Events[:t.maxEvents]
	}
	// Truncate oversized data fields.
	for i := range s.Events {
		d := s.Events[i].Data
		if len(d) > t.maxBytes {
			// Trim to a valid UTF-8 boundary within the budget.
			budget := t.maxBytes - len(truncationMarker)
			if budget < 0 {
				budget = 0
			}
			trimmed := d[:budget]
			// Walk back to a valid rune boundary.
			for !utf8.ValidString(trimmed) && len(trimmed) > 0 {
				trimmed = trimmed[:len(trimmed)-1]
			}
			s.Events[i].Data = trimmed + truncationMarker
		}
	}
}
