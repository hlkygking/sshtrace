// Package rotate provides session log rotation based on age or count limits.
package rotate

import (
	"errors"
	"time"

	"sshtrace/internal/session"
)

// Policy defines when sessions should be rotated out of storage.
type Policy struct {
	// MaxAge is the maximum age of a session before it is eligible for rotation.
	MaxAge time.Duration
	// MaxCount is the maximum number of sessions to retain (0 means unlimited).
	MaxCount int
}

// DefaultPolicy returns a sensible default rotation policy.
func DefaultPolicy() Policy {
	return Policy{
		MaxAge:   30 * 24 * time.Hour, // 30 days
		MaxCount: 10000,
	}
}

// Rotator applies a rotation policy to a slice of sessions.
type Rotator struct {
	policy Policy
}

// New creates a new Rotator with the given policy.
func New(p Policy) (*Rotator, error) {
	if p.MaxAge < 0 {
		return nil, errors.New("rotate: MaxAge must be non-negative")
	}
	if p.MaxCount < 0 {
		return nil, errors.New("rotate: MaxCount must be non-negative")
	}
	return &Rotator{policy: p}, nil
}

// Apply returns the sessions that should be kept and those that should be removed.
// Sessions are assumed to be ordered oldest-first.
func (r *Rotator) Apply(sessions []*session.Session) (keep, remove []*session.Session) {
	now := time.Now()

	// First pass: filter by age.
	for _, s := range sessions {
		if r.policy.MaxAge > 0 && now.Sub(s.StartedAt) > r.policy.MaxAge {
			remove = append(remove, s)
		} else {
			keep = append(keep, s)
		}
	}

	// Second pass: enforce MaxCount on the kept slice (remove oldest first).
	if r.policy.MaxCount > 0 && len(keep) > r.policy.MaxCount {
		excess := keep[:len(keep)-r.policy.MaxCount]
		remove = append(remove, excess...)
		keep = keep[len(keep)-r.policy.MaxCount:]
	}

	return keep, remove
}
