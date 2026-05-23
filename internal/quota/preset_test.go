package quota_test

import (
	"testing"
	"time"

	"sshtrace/internal/quota"
)

func TestPresetStrict(t *testing.T) {
	l := quota.PresetLimits(quota.PresetStrict)
	if l.MaxSessions != 10 {
		t.Fatalf("strict MaxSessions: want 10, got %d", l.MaxSessions)
	}
	if l.MaxEventsPerSession != 500 {
		t.Fatalf("strict MaxEventsPerSession: want 500, got %d", l.MaxEventsPerSession)
	}
	if l.MaxSessionDuration != time.Hour {
		t.Fatalf("strict MaxSessionDuration: want 1h, got %v", l.MaxSessionDuration)
	}
}

func TestPresetRelaxed(t *testing.T) {
	l := quota.PresetLimits(quota.PresetRelaxed)
	if l.MaxSessions != 500 {
		t.Fatalf("relaxed MaxSessions: want 500, got %d", l.MaxSessions)
	}
	if l.MaxSessionDuration != 24*time.Hour {
		t.Fatalf("relaxed MaxSessionDuration: want 24h, got %v", l.MaxSessionDuration)
	}
}

func TestPresetUnlimited(t *testing.T) {
	l := quota.PresetLimits(quota.PresetUnlimited)
	if l.MaxSessions != 0 || l.MaxEventsPerSession != 0 || l.MaxSessionDuration != 0 {
		t.Fatal("unlimited preset should have all zero limits")
	}
}

func TestPresetUnknownFallsBackToDefault(t *testing.T) {
	got := quota.PresetLimits("nonexistent")
	want := quota.DefaultLimits()
	if got.MaxSessions != want.MaxSessions {
		t.Fatalf("unknown preset MaxSessions: want %d, got %d", want.MaxSessions, got.MaxSessions)
	}
	if got.MaxEventsPerSession != want.MaxEventsPerSession {
		t.Fatalf("unknown preset MaxEventsPerSession: want %d, got %d", want.MaxEventsPerSession, got.MaxEventsPerSession)
	}
}

func TestPresetEnforcerIntegration(t *testing.T) {
	e := quota.New(quota.PresetLimits(quota.PresetStrict))
	for i := 0; i < 10; i++ {
		e.RecordSession("grace")
	}
	if err := e.CheckSession("grace"); err != quota.ErrQuotaExceeded {
		t.Fatalf("expected quota exceeded after 10 strict sessions, got %v", err)
	}
}
