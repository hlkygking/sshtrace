// Package classify provides risk-level classification for SSH session events.
//
// Each event's data is matched against curated regular-expression lists to
// produce one of three levels:
//
//   - LevelSafe     – routine, read-only activity (ls, cat, echo, …)
//   - LevelModerate – commands that modify state or install software
//   - LevelDangerous – destructive or privilege-escalating operations
//
// Usage:
//
//	level := classify.Classify("sudo rm -rf /tmp/work")
//	// level == classify.LevelDangerous
//
//	results := classify.Analyse(sess)
//	for _, r := range results {
//		fmt.Printf("event %d: %s\n", r.EventIndex, r.Level)
//	}
package classify
