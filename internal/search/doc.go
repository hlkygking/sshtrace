// Package search provides full-text search capabilities over recorded SSH
// session events.
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
package search
