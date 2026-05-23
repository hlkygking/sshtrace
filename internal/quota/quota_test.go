package quota_test

import (
	"testing"
	"time"

	"sshtrace/internal/quota"
)

func TestDefaultLimits(t *testing.T) {
	l := quota.DefaultLimits()
	if l.MaxSessions <= 0 {
		t.Fatal("expected positive MaxSessions")
	}
	if l.MaxEventsPerSession <= 0 {
		t.Fatal("expected positive MaxEventsPerSession")
	}
	if l.MaxSessionDuration <= 0 {
		t.Fatal("expected positive MaxSessionDuration")
	}
}

func TestCheckSessionAllowed(t *testing.T) {
	e := quota.New(quota.Limits{MaxSessions: 5})
	if err := e.CheckSession("alice"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheckSessionExceeded(t *testing.T) {
	e := quota.New(quota.Limits{MaxSessions: 2})
	e.RecordSession("bob")
	e.RecordSession("bob")
	if err := e.CheckSession("bob"); err != quota.ErrQuotaExceeded {
		t.Fatalf("expected ErrQuotaExceeded, got %v", err)
	}
}

func TestRecordSessionIncrementsCount(t *testing.T) {
	e := quota.New(quota.DefaultLimits())
	e.RecordSession("carol")
	e.RecordSession("carol")
	if got := e.Sessions("carol"); got != 2 {
		t.Fatalf("expected 2 sessions, got %d", got)
	}
}

func TestCheckEventsAllowed(t *testing.T) {
	e := quota.New(quota.Limits{MaxEventsPerSession: 100})
	if err := e.CheckEvents(50); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheckEventsExceeded(t *testing.T) {
	e := quota.New(quota.Limits{MaxEventsPerSession: 10})
	if err := e.CheckEvents(11); err != quota.ErrQuotaExceeded {
		t.Fatalf("expected ErrQuotaExceeded, got %v", err)
	}
}

func TestCheckDurationAllowed(t *testing.T) {
	e := quota.New(quota.Limits{MaxSessionDuration: time.Hour})
	if err := e.CheckDuration(30 * time.Minute); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheckDurationExceeded(t *testing.T) {
	e := quota.New(quota.Limits{MaxSessionDuration: time.Hour})
	if err := e.CheckDuration(2 * time.Hour); err != quota.ErrQuotaExceeded {
		t.Fatalf("expected ErrQuotaExceeded, got %v", err)
	}
}

func TestIsolatedUsersDoNotShareCounts(t *testing.T) {
	e := quota.New(quota.Limits{MaxSessions: 3})
	e.RecordSession("dave")
	e.RecordSession("dave")
	if err := e.CheckSession("eve"); err != nil {
		t.Fatalf("eve should not be affected by dave's usage: %v", err)
	}
}

func TestZeroLimitMeansUnlimited(t *testing.T) {
	e := quota.New(quota.Limits{MaxSessions: 0})
	for i := 0; i < 1000; i++ {
		e.RecordSession("frank")
	}
	if err := e.CheckSession("frank"); err != nil {
		t.Fatalf("zero limit should be unlimited, got %v", err)
	}
}
