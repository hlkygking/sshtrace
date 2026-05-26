package window

import "time"

// PresetAggregator returns a pre-configured Aggregator for common window sizes.
// Supported names: "second", "minute", "hour". Unknown names fall back to
// a one-minute window.
func PresetAggregator(name string) *Aggregator {
	sizes := map[string]time.Duration{
		"second": time.Second,
		"minute": time.Minute,
		"hour":   time.Hour,
	}
	d, ok := sizes[name]
	if !ok {
		d = time.Minute
	}
	// size is always positive here, so we can ignore the error.
	agg, _ := New(d)
	return agg
}
