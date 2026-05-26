// Package window implements sliding-window aggregation for SSH session events.
//
// Use New to create an Aggregator with a fixed window duration, then call
// Aggregate to partition a session's events into consecutive time buckets.
// Each bucket exposes the event count, unique users, and raw command strings
// observed within that interval.
//
// Example:
//
//	agg, err := window.New(30 * time.Second)
//	if err != nil { ... }
//
//	results, err := agg.Aggregate(sess)
//	for _, r := range results {
//		fmt.Printf("%s – %s : %d events\n", r.Start, r.End, r.EventCount)
//	}
package window
