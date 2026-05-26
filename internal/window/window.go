// Package window provides sliding-window aggregation over SSH session events.
// It groups events that fall within a configurable time window and exposes
// aggregate statistics such as event count, unique commands, and active users.
package window

import (
	"errors"
	"time"

	"sshtrace/internal/session"
)

// Result holds the aggregated data for a single time window.
type Result struct {
	Start       time.Time
	End         time.Time
	EventCount  int
	UniqueUsers map[string]struct{}
	Commands    []string
}

// Aggregator groups session events into fixed-size time windows.
type Aggregator struct {
	size time.Duration
}

// New creates an Aggregator with the given window size.
// Returns an error if size is zero or negative.
func New(size time.Duration) (*Aggregator, error) {
	if size <= 0 {
		return nil, errors.New("window: size must be positive")
	}
	return &Aggregator{size: size}, nil
}

// Aggregate partitions the events of s into non-overlapping windows of a.size
// and returns one Result per window that contains at least one event.
func (a *Aggregator) Aggregate(s *session.Session) ([]Result, error) {
	if s == nil {
		return nil, errors.New("window: session must not be nil")
	}

	if len(s.Events) == 0 {
		return nil, nil
	}

	buckets := make(map[int64]*Result)
	var order []int64

	for _, ev := range s.Events {
		key := ev.Timestamp.UnixNano() / int64(a.size)
		if _, exists := buckets[key]; !exists {
			start := time.Unix(0, key*int64(a.size))
			buckets[key] = &Result{
				Start:       start,
				End:         start.Add(a.size),
				UniqueUsers: make(map[string]struct{}),
			}
			order = append(order, key)
		}
		b := buckets[key]
		b.EventCount++
		b.UniqueUsers[s.User] = struct{}{}
		if ev.Data != "" {
			b.Commands = append(b.Commands, ev.Data)
		}
	}

	results := make([]Result, 0, len(order))
	for _, k := range order {
		results = append(results, *buckets[k])
	}
	return results, nil
}
