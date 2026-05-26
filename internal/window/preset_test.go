package window_test

import (
	"testing"
	"time"

	"sshtrace/internal/session"
	"sshtrace/internal/window"
)

func TestPresetSecondBuckets(t *testing.T) {
	agg := window.PresetAggregator("second")
	base := time.Now().Truncate(time.Second)
	s := &session.Session{
		User: "dave",
		Events: []session.Event{
			{Type: "output", Data: "a", Timestamp: base},
			{Type: "output", Data: "b", Timestamp: base.Add(2 * time.Second)},
		},
	}
	results, err := agg.Aggregate(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 second-buckets, got %d", len(results))
	}
}

func TestPresetMinuteSingleBucket(t *testing.T) {
	agg := window.PresetAggregator("minute")
	base := time.Now().Truncate(time.Minute)
	s := &session.Session{
		User: "eve",
		Events: []session.Event{
			{Type: "output", Data: "x", Timestamp: base.Add(5 * time.Second)},
			{Type: "output", Data: "y", Timestamp: base.Add(30 * time.Second)},
		},
	}
	results, err := agg.Aggregate(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 minute-bucket, got %d", len(results))
	}
	if results[0].EventCount != 2 {
		t.Errorf("expected 2 events, got %d", results[0].EventCount)
	}
}
