// Package rotate implements log rotation for SSH session records.
//
// A Policy specifies two independent limits:
//
//   - MaxAge: sessions older than this duration are eligible for removal.
//   - MaxCount: if more sessions than this limit are present the oldest
//     ones are removed until the count is within the limit.
//
// Either limit can be disabled by setting it to zero.
//
// Example usage:
//
//	pol := rotate.DefaultPolicy()
//	rot, _ := rotate.New(pol)
//	keep, remove := rot.Apply(allSessions)
package rotate
