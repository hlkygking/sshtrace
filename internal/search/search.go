// Package search provides full-text search across SSH session events.
package search

import (
	"strings"

	"sshtrace/internal/session"
)

// Options configures a search query.
type Options struct {
	// Query is the substring to search for in event data.
	Query string
	// CaseSensitive controls whether matching is case-sensitive.
	CaseSensitive bool
	// EventType filters results to a specific event type ("input", "output", or "" for all).
	EventType string
}

// Result holds a matched event along with its parent session.
type Result struct {
	Session *session.Session
	Event   session.Event
}

// Search scans the provided sessions for events matching opts.
// It returns a slice of Result, one per matching event.
func Search(sessions []*session.Session, opts Options) []Result {
	var results []Result

	query := opts.Query
	if !opts.CaseSensitive {
		query = strings.ToLower(query)
	}

	for _, s := range sessions {
		for _, ev := range s.Events {
			if opts.EventType != "" && ev.Type != opts.EventType {
				continue
			}

			data := ev.Data
			if !opts.CaseSensitive {
				data = strings.ToLower(data)
			}

			if strings.Contains(data, query) {
				results = append(results, Result{
					Session: s,
					Event:   ev,
				})
			}
		}
	}

	return results
}
