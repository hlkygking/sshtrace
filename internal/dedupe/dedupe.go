// Package dedupe provides deduplication of SSH session events based on
// content hashing, removing consecutive identical events within a session.
package dedupe

import (
	"crypto/sha256"
	"fmt"

	"sshtrace/internal/session"
)

// Deduplicator removes consecutive duplicate events from a session.
type Deduplicator struct {
	// WindowSize is the number of recent event hashes to track.
	// A value of 1 removes only consecutive duplicates.
	WindowSize int
}

// New returns a Deduplicator with the given window size.
// windowSize must be >= 1.
func New(windowSize int) (*Deduplicator, error) {
	if windowSize < 1 {
		return nil, fmt.Errorf("dedupe: windowSize must be >= 1, got %d", windowSize)
	}
	return &Deduplicator{WindowSize: windowSize}, nil
}

// Apply removes duplicate events from the session in-place and returns
// the modified session. The original event slice is replaced.
func (d *Deduplicator) Apply(s *session.Session) (*session.Session, error) {
	if s == nil {
		return nil, fmt.Errorf("dedupe: session must not be nil")
	}

	window := make([]string, 0, d.WindowSize)
	filtered := make([]session.Event, 0, len(s.Events))

	for _, ev := range s.Events {
		h := hashEvent(ev)
		if !inWindow(window, h) {
			filtered = append(filtered, ev)
			window = appendWindow(window, h, d.WindowSize)
		}
	}

	s.Events = filtered
	return s, nil
}

func hashEvent(ev session.Event) string {
	raw := fmt.Sprintf("%s|%s", ev.Type, ev.Data)
	sum := sha256.Sum256([]byte(raw))
	return fmt.Sprintf("%x", sum)
}

func inWindow(window []string, h string) bool {
	for _, w := range window {
		if w == h {
			return true
		}
	}
	return false
}

func appendWindow(window []string, h string, max int) []string {
	window = append(window, h)
	if len(window) > max {
		window = window[len(window)-max:]
	}
	return window
}
