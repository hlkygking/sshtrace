// Package search provides full-text search capabilities over recorded SSH
// session events.
//
// The primary entry point is the Search function, which accepts a slice of
// sessions and an Options struct to control matching behavior. Searches can
// be scoped to a specific event type (e.g. "input", "output") or run across
// all event types when EventType is left empty.
//
// Usage:
//
//	results := search.Search(sessions, search.Options{
//		Query:         "sudo",
//		CaseSensitive: false,
//		EventType:     "input",
//	})
//	for _, r := range results {
//		fmt.Printf("[%s] %s: %s\n", r.Session.ID, r.Event.Type, r.Event.Data)
//	}
//
// Search returns one Result per matching event, each carrying a reference to
// the parent Session so callers can access session metadata alongside the
// matched content.
//
// Result ordering mirrors the order in which sessions and their events were
// provided; no additional sorting is applied.
package search
