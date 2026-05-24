package ratelimit_test

import (
	"testing"

	"github.com/sshtrace/sshtrace/internal/ratelimit"
)

func TestPresetStrict(t *testing.T) {
	cfg := ratelimit.PresetConfig("strict")
	if cfg.MaxConnections >= 10 {
		t.Errorf("strict preset should have low MaxConnections, got %d", cfg.MaxConnections)
	}
}

func TestPresetRelaxed(t *testing.T) {
	cfg := ratelimit.PresetConfig("relaxed")
	if cfg.MaxConnections <= 10 {
		t.Errorf("relaxed preset should have high MaxConnections, got %d", cfg.MaxConnections)
	}
}

func TestPresetUnlimited(t *testing.T) {
	cfg := ratelimit.PresetConfig("unlimited")
	l, err := ratelimit.New(cfg)
	if err != nil {
		t.Fatalf("unlimited preset should be valid: %v", err)
	}
	for i := 0; i < 100; i++ {
		if err := l.Allow("user"); err != nil {
			t.Fatalf("unlimited should never rate limit, failed at %d: %v", i, err)
		}
	}
}

func TestPresetUnknownFallsBackToDefault(t *testing.T) {
	cfg := ratelimit.PresetConfig("unknown-preset")
	def := ratelimit.DefaultConfig()
	if cfg.MaxConnections != def.MaxConnections || cfg.Window != def.Window {
		t.Errorf("unknown preset should equal default, got %+v", cfg)
	}
}
