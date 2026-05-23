// Package quota enforces per-user session and event storage limits.
package quota

import (
	"errors"
	"sync"
	"time"
)

// ErrQuotaExceeded is returned when a user exceeds their allowed quota.
var ErrQuotaExceeded = errors.New("quota exceeded")

// Limits defines the maximum allowed values for a user.
type Limits struct {
	MaxSessions  int
	MaxEventsPerSession int
	MaxSessionDuration  time.Duration
}

// DefaultLimits returns sensible default quota limits.
func DefaultLimits() Limits {
	return Limits{
		MaxSessions:         100,
		MaxEventsPerSession: 10000,
		MaxSessionDuration:  8 * time.Hour,
	}
}

// Enforcer tracks usage and enforces quota limits per user.
type Enforcer struct {
	mu     sync.Mutex
	limits Limits
	usage  map[string]*usage
}

type usage struct {
	sessions int
}

// New creates a new Enforcer with the given limits.
func New(limits Limits) *Enforcer {
	return &Enforcer{
		limits: limits,
		usage:  make(map[string]*usage),
	}
}

// CheckSession returns ErrQuotaExceeded if the user has reached their session limit.
func (e *Enforcer) CheckSession(user string) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	u := e.getOrCreate(user)
	if e.limits.MaxSessions > 0 && u.sessions >= e.limits.MaxSessions {
		return ErrQuotaExceeded
	}
	return nil
}

// RecordSession increments the session count for the given user.
func (e *Enforcer) RecordSession(user string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.getOrCreate(user).sessions++
}

// CheckEvents returns ErrQuotaExceeded if count exceeds the per-session event limit.
func (e *Enforcer) CheckEvents(count int) error {
	if e.limits.MaxEventsPerSession > 0 && count > e.limits.MaxEventsPerSession {
		return ErrQuotaExceeded
	}
	return nil
}

// CheckDuration returns ErrQuotaExceeded if d exceeds the max session duration.
func (e *Enforcer) CheckDuration(d time.Duration) error {
	if e.limits.MaxSessionDuration > 0 && d > e.limits.MaxSessionDuration {
		return ErrQuotaExceeded
	}
	return nil
}

// Sessions returns the current session count for the user.
func (e *Enforcer) Sessions(user string) int {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.getOrCreate(user).sessions
}

func (e *Enforcer) getOrCreate(user string) *usage {
	if _, ok := e.usage[user]; !ok {
		e.usage[user] = &usage{}
	}
	return e.usage[user]
}
