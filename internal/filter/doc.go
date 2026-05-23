// Package filter provides utilities for querying and narrowing collections
// of SSH sessions based on common audit criteria such as user name, remote IP
// address, time range, and minimum event count.
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
