// Package merge provides utilities for combining multiple SSH sessions
// into a single unified session, ordered by event timestamp.
package merge

import (
	"errors"
	"sort"
	"time"

	"sshtrace/internal/session"
)

// Options controls how sessions are merged.
type Options struct {
	// PreserveIDs keeps original session IDs in event metadata.
	PreserveIDs bool
	// DeduplicateWindow discards events with identical data within this
	// duration of each other. Zero disables deduplication.
	DeduplicateWindow time.Duration
}

// DefaultOptions returns sensible merge defaults.
func DefaultOptions() Options {
	return Options{
		PreserveIDs:       true,
		DeduplicateWindow: 0,
	}
}

// Merger combines sessions.
type Merger struct {
	opts Options
}

// New creates a Merger with the given options.
func New(opts Options) (*Merger, error) {
	if opts.DeduplicateWindow < 0 {
		return nil, errors.New("merge: DeduplicateWindow must not be negative")
	}
	return &Merger{opts: opts}, nil
}

// Merge combines all provided sessions into one, sorting events by timestamp.
// The resulting session inherits the user and remote IP from the first session.
func (m *Merger) Merge(sessions []*session.Session) (*session.Session, error) {
	if len(sessions) == 0 {
		return nil, errors.New("merge: no sessions provided")
	}
	for i, s := range sessions {
		if s == nil {
			return nil, errors.New("merge: nil session at index " + itoa(i))
		}
	}

	base := sessions[0]
	out := session.New(base.User, base.RemoteIP)

	type tagged struct {
		event  session.Event
		origin string
	}

	var all []tagged
	for _, s := range sessions {
		for _, e := range s.Events {
			all = append(all, tagged{event: e, origin: s.ID})
		}
	}

	sort.SliceStable(all, func(i, j int) bool {
		return all[i].event.Timestamp.Before(all[j].event.Timestamp)
	})

	seen := map[string]time.Time{}
	for _, t := range all {
		e := t.event
		if m.opts.PreserveIDs && e.Meta == nil {
			e.Meta = map[string]string{}
		}
		if m.opts.PreserveIDs {
			e.Meta["origin_session"] = t.origin
		}
		if m.opts.DeduplicateWindow > 0 {
			key := e.Kind + ":" + e.Data
			if last, ok := seen[key]; ok && e.Timestamp.Sub(last) < m.opts.DeduplicateWindow {
				continue
			}
			seen[key] = e.Timestamp
		}
		out.AddEvent(e)
	}

	return out, nil
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	b := make([]byte, 0, 4)
	for n > 0 {
		b = append([]byte{byte('0' + n%10)}, b...)
		n /= 10
	}
	return string(b)
}
