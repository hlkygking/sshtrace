// Package filter provides utilities for querying and narrowing collections
// of SSH sessions based on common audit criteria such as user name, remote IP
// address, time range, and minimum event count.
//
// The primary entry point is [Filter], which accepts a slice of sessions and
// a [Criteria] struct. Only sessions satisfying all non-zero criteria fields
// are returned.
//
// Criteria fields:
//
//   - UserName:  exact match against the session's authenticated user
//   - RemoteIP:  exact match against the client's remote IP address
//   - Since:     only sessions that started at or after this time are included
//   - MinEvents: only sessions with at least this many recorded events are included
//
// Example usage:
//
//	sessions, _ := storage.ListSince(store, time.Time{})
//	matches := filter.Filter(sessions, filter.Criteria{
//		UserName:  "alice",
//		MinEvents: 1,
//		Since:     time.Now().Add(-24 * time.Hour),
//	})
package filter
