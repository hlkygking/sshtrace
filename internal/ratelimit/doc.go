// Package ratelimit implements a sliding-window rate limiter for SSH connections.
//
// It can be keyed by any string identifier — typically a username or remote IP
// address — and enforces a maximum number of connections within a configurable
// time window.
//
// Example usage:
//
//	limiter, err := ratelimit.New(ratelimit.DefaultConfig())
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	if err := limiter.Allow(remoteIP); err != nil {
//		// reject connection
//	}
package ratelimit
