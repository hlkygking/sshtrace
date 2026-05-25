// Package correlate links related SSH sessions together based on shared
// attributes such as user, source IP, or overlapping time windows.
package correlate

import (
	"errors"
	"time"

	"sshtrace/internal/session"
)

// Strategy controls how sessions are grouped into correlation clusters.
type Strategy int

const (
	// ByUser groups sessions that share the same username.
	ByUser Strategy = iota
	// ByIP groups sessions that share the same source IP address.
	ByIP
	// ByTimeWindow groups sessions whose active periods overlap within a
	// configurable tolerance.
	ByTimeWindow
)

// Options configures the Correlator.
type Options struct {
	Strategy  Strategy
	// Window is only used when Strategy == ByTimeWindow.
	Window    time.Duration
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Strategy: ByUser,
		Window:   5 * time.Minute,
	}
}

// Correlator groups sessions into clusters.
type Correlator struct {
	opts Options
}

// New creates a Correlator with the given options.
func New(opts Options) (*Correlator, error) {
	if opts.Strategy == ByTimeWindow && opts.Window <= 0 {
		return nil, errors.New("correlate: window must be positive for ByTimeWindow strategy")
	}
	return &Correlator{opts: opts}, nil
}

// Group partitions sessions into clusters. Each inner slice is one cluster.
func (c *Correlator) Group(sessions []*session.Session) [][]*session.Session {
	if len(sessions) == 0 {
		return nil
	}
	switch c.opts.Strategy {
	case ByIP:
		return groupBy(sessions, func(s *session.Session) string { return s.SourceIP })
	case ByTimeWindow:
		return groupByTimeWindow(sessions, c.opts.Window)
	default: // ByUser
		return groupBy(sessions, func(s *session.Session) string { return s.User })
	}
}

func groupBy(sessions []*session.Session, key func(*session.Session) string) [][]*session.Session {
	order := []string{}
	buckets := map[string][]*session.Session{}
	for _, s := range sessions {
		k := key(s)
		if _, exists := buckets[k]; !exists {
			order = append(order, k)
		}
		buckets[k] = append(buckets[k], s)
	}
	out := make([][]*session.Session, 0, len(order))
	for _, k := range order {
		out = append(out, buckets[k])
	}
	return out
}

func groupByTimeWindow(sessions []*session.Session, window time.Duration) [][]*session.Session {
	used := make([]bool, len(sessions))
	var clusters [][]*session.Session
	for i, s := range sessions {
		if used[i] {
			continue
		}
		cluster := []*session.Session{s}
		used[i] = true
		for j := i + 1; j < len(sessions); j++ {
			if used[j] {
				continue
			}
			if overlaps(s, sessions[j], window) {
				cluster = append(cluster, sessions[j])
				used[j] = true
			}
		}
		clusters = append(clusters, cluster)
	}
	return clusters
}

func overlaps(a, b *session.Session, window time.Duration) bool {
	aStart := a.StartedAt
	aEnd := a.StartedAt.Add(window)
	if !a.EndedAt.IsZero() {
		aEnd = a.EndedAt.Add(window)
	}
	bStart := b.StartedAt
	return !bStart.After(aEnd) && !aStart.After(bStart.Add(window))
}
