package ratelimit_test

import (
	"testing"
	"time"

	"github.com/sshtrace/sshtrace/internal/ratelimit"
)

func TestDefaultConfig(t *testing.T) {
	cfg := ratelimit.DefaultConfig()
	if cfg.MaxConnections <= 0 {
		t.Errorf("expected positive MaxConnections, got %d", cfg.MaxConnections)
	}
	if cfg.Window <= 0 {
		t.Errorf("expected positive Window, got %v", cfg.Window)
	}
}

func TestNewInvalidMaxConnections(t *testing.T) {
	_, err := ratelimit.New(ratelimit.Config{MaxConnections: 0, Window: time.Minute})
	if err == nil {
		t.Error("expected error for zero MaxConnections")
	}
}

func TestNewInvalidWindow(t *testing.T) {
	_, err := ratelimit.New(ratelimit.Config{MaxConnections: 5, Window: 0})
	if err == nil {
		t.Error("expected error for zero Window")
	}
}

func TestAllowUnderLimit(t *testing.T) {
	l, err := ratelimit.New(ratelimit.Config{MaxConnections: 3, Window: time.Minute})
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 3; i++ {
		if err := l.Allow("user1"); err != nil {
			t.Errorf("expected nil on attempt %d, got %v", i+1, err)
		}
	}
}

func TestAllowExceedsLimit(t *testing.T) {
	l, err := ratelimit.New(ratelimit.Config{MaxConnections: 2, Window: time.Minute})
	if err != nil {
		t.Fatal(err)
	}
	_ = l.Allow("alice")
	_ = l.Allow("alice")
	if err := l.Allow("alice"); err == nil {
		t.Error("expected ErrRateLimited on third attempt")
	}
}

func TestAllowIndependentKeys(t *testing.T) {
	l, _ := ratelimit.New(ratelimit.Config{MaxConnections: 1, Window: time.Minute})
	if err := l.Allow("alice"); err != nil {
		t.Errorf("alice should be allowed: %v", err)
	}
	if err := l.Allow("bob"); err != nil {
		t.Errorf("bob should be allowed: %v", err)
	}
}

func TestReset(t *testing.T) {
	l, _ := ratelimit.New(ratelimit.Config{MaxConnections: 1, Window: time.Minute})
	_ = l.Allow("alice")
	l.Reset("alice")
	if err := l.Allow("alice"); err != nil {
		t.Errorf("expected nil after reset, got %v", err)
	}
}

func TestCount(t *testing.T) {
	l, _ := ratelimit.New(ratelimit.Config{MaxConnections: 10, Window: time.Minute})
	_ = l.Allow("user")
	_ = l.Allow("user")
	if c := l.Count("user"); c != 2 {
		t.Errorf("expected count 2, got %d", c)
	}
}
