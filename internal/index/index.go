// Package index provides fast lookup of sessions by user, IP, or time range
// using in-memory inverted indexes built from a session slice.
package index

import (
	"strings"
	"time"

	"sshtrace/internal/session"
)

// Index holds pre-built lookup structures for a collection of sessions.
type Index struct {
	byUser map[string][]*session.Session
	byIP   map[string][]*session.Session
	all    []*session.Session
}

// Build constructs an Index from the provided sessions.
func Build(sessions []*session.Session) *Index {
	idx := &Index{
		byUser: make(map[string][]*session.Session),
		byIP:   make(map[string][]*session.Session),
		all:    sessions,
	}
	for _, s := range sessions {
		key := strings.ToLower(s.User)
		idx.byUser[key] = append(idx.byUser[key], s)
		idx.byIP[s.ClientIP] = append(idx.byIP[s.ClientIP], s)
	}
	return idx
}

// LookupUser returns all sessions for the given username (case-insensitive).
func (idx *Index) LookupUser(user string) []*session.Session {
	return idx.byUser[strings.ToLower(user)]
}

// LookupIP returns all sessions originating from the given IP address.
func (idx *Index) LookupIP(ip string) []*session.Session {
	return idx.byIP[ip]
}

// LookupRange returns all sessions whose start time falls within [from, to].
func (idx *Index) LookupRange(from, to time.Time) []*session.Session {
	var result []*session.Session
	for _, s := range idx.all {
		if !s.StartedAt.Before(from) && !s.StartedAt.After(to) {
			result = append(result, s)
		}
	}
	return result
}

// Size returns the total number of indexed sessions.
func (idx *Index) Size() int {
	return len(idx.all)
}
