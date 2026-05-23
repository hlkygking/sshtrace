// Package diff compares two SSH session event streams and reports
// differences as a list of deltas.
//
// Usage:
//
//	res := diff.Compare(sessionA, sessionB)
//	if !res.Equal() {
//		fmt.Println(res.Summary())
//	}
//
// Each Delta carries the index of the differing event, the change type
// (added, removed, or changed), and pointers to the original events from
// each session for further inspection.
package diff
