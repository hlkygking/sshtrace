// Package ratelimit provides per-user and per-IP rate limiting for SSH sessions.
package ratelimit

import (
	"errors"
	"sync"
	"time"
)

// ErrRateLimited is returned when a connection exceeds the allowed rate.
var ErrRateLimited = errors.New("rate limit exceeded")

// Config holds rate limiting parameters.
type Config struct {
	// MaxConnections is the maximum number of connections allowed per window.
	MaxConnections int
	// Window is the duration of the sliding window.
	Window time.Duration
}

// DefaultConfig returns a sensible default rate limit configuration.
func DefaultConfig() Config {
	return Config{
		MaxConnections: 10,
		Window:         time.Minute,
	}
}

// Limiter tracks connection attempts per key (user or IP).
type Limiter struct {
	cfg    Config
	mu     sync.Mutex
	bucket map[string][]time.Time
}

// New creates a new Limiter with the given configuration.
func New(cfg Config) (*Limiter, error) {
	if cfg.MaxConnections <= 0 {
		return nil, errors.New("MaxConnections must be greater than zero")
	}
	if cfg.Window <= 0 {
		return nil, errors.New("Window must be greater than zero")
	}
	return &Limiter{
		cfg:    cfg,
		bucket: make(map[string][]time.Time),
	}, nil
}

// Allow checks whether a new connection from key is permitted.
// It records the attempt and returns ErrRateLimited if the limit is exceeded.
func (l *Limiter) Allow(key string) error {
	now := time.Now()
	cutoff := now.Add(-l.cfg.Window)

	l.mu.Lock()
	defer l.mu.Unlock()

	times := l.bucket[key]
	filtered := times[:0]
	for _, t := range times {
		if t.After(cutoff) {
			filtered = append(filtered, t)
		}
	}

	if len(filtered) >= l.cfg.MaxConnections {
		l.bucket[key] = filtered
		return ErrRateLimited
	}

	l.bucket[key] = append(filtered, now)
	return nil
}

// Reset clears the rate limit state for a given key.
func (l *Limiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.bucket, key)
}

// Count returns the number of recorded attempts within the window for key.
func (l *Limiter) Count(key string) int {
	cutoff := time.Now().Add(-l.cfg.Window)
	l.mu.Lock()
	defer l.mu.Unlock()
	count := 0
	for _, t := range l.bucket[key] {
		if t.After(cutoff) {
			count++
		}
	}
	return count
}
