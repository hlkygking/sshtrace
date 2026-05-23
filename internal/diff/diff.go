// Package diff provides utilities for comparing two SSH sessions
// and highlighting differences in their event streams.
package diff

import (
	"fmt"
	"strings"

	"sshtrace/internal/session"
)

// ChangeType indicates the kind of difference found.
type ChangeType string

const (
	Added   ChangeType = "added"
	Removed ChangeType = "removed"
	Changed ChangeType = "changed"
)

// Delta represents a single difference between two sessions.
type Delta struct {
	Index  int
	Type   ChangeType
	Left   *session.Event
	Right  *session.Event
}

// Result holds the full diff output between two sessions.
type Result struct {
	LeftID  string
	RightID string
	Deltas  []Delta
}

// Equal reports whether the diff found no differences.
func (r *Result) Equal() bool {
	return len(r.Deltas) == 0
}

// Summary returns a human-readable summary of the diff.
func (r *Result) Summary() string {
	if r.Equal() {
		return fmt.Sprintf("sessions %s and %s are identical", r.LeftID, r.RightID)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "diff %s..%s: %d delta(s)\n", r.LeftID, r.RightID, len(r.Deltas))
	for _, d := range r.Deltas {
		switch d.Type {
		case Added:
			fmt.Fprintf(&sb, "  [%d] + %s\n", d.Index, d.Right.Data)
		case Removed:
			fmt.Fprintf(&sb, "  [%d] - %s\n", d.Index, d.Left.Data)
		case Changed:
			fmt.Fprintf(&sb, "  [%d] ~ %s -> %s\n", d.Index, d.Left.Data, d.Right.Data)
		}
	}
	return strings.TrimRight(sb.String(), "\n")
}

// Compare diffs the event streams of two sessions.
// It aligns events by index and reports additions, removals, and changes.
func Compare(a, b *session.Session) *Result {
	res := &Result{
		LeftID:  a.ID,
		RightID: b.ID,
	}

	la, lb := len(a.Events), len(b.Events)
	max := la
	if lb > max {
		max = lb
	}

	for i := 0; i < max; i++ {
		switch {
		case i >= la:
			res.Deltas = append(res.Deltas, Delta{Index: i, Type: Added, Right: &b.Events[i]})
		case i >= lb:
			res.Deltas = append(res.Deltas, Delta{Index: i, Type: Removed, Left: &a.Events[i]})
		case a.Events[i].Data != b.Events[i].Data:
			res.Deltas = append(res.Deltas, Delta{Index: i, Type: Changed, Left: &a.Events[i], Right: &b.Events[i]})
		}
	}
	return res
}
