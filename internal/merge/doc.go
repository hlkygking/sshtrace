// Package merge provides a Merger that combines multiple SSH session records
// into a single unified session.
//
// Events from all input sessions are interleaved in chronological order by
// their timestamps. The merged session inherits the user and remote IP from
// the first provided session.
//
// Optional deduplication removes events whose Kind and Data are identical
// within a configurable time window, which is useful when the same command
// appears in overlapping capture streams.
//
// Example usage:
//
//	m, _ := merge.New(merge.DefaultOptions())
//	result, err := m.Merge([]*session.Session{s1, s2})
package merge
