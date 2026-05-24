package schedule_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/sshtrace/sshtrace/internal/schedule"
)

func TestAddEmptyNameReturnsError(t *testing.T) {
	s := schedule.New()
	err := s.Add("", time.Second, func() error { return nil })
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestAddZeroIntervalReturnsError(t *testing.T) {
	s := schedule.New()
	err := s.Add("job", 0, func() error { return nil })
	if err == nil {
		t.Fatal("expected error for zero interval")
	}
}

func TestAddNilTaskReturnsError(t *testing.T) {
	s := schedule.New()
	err := s.Add("job", time.Second, nil)
	if err == nil {
		t.Fatal("expected error for nil task")
	}
}

func TestAddRegistersJob(t *testing.T) {
	s := schedule.New()
	if err := s.Add("cleanup", time.Minute, func() error { return nil }); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Jobs()) != 1 {
		t.Fatalf("expected 1 job, got %d", len(s.Jobs()))
	}
	if s.Jobs()[0].Name != "cleanup" {
		t.Errorf("expected job name 'cleanup', got %q", s.Jobs()[0].Name)
	}
}

func TestJobExecutesOnInterval(t *testing.T) {
	s := schedule.New()
	var count int64
	err := s.Add("tick", 20*time.Millisecond, func() error {
		atomic.AddInt64(&count, 1)
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s.Start()
	time.Sleep(70 * time.Millisecond)
	s.Stop()

	got := atomic.LoadInt64(&count)
	if got < 2 {
		t.Errorf("expected at least 2 executions, got %d", got)
	}
}

func TestStopHaltsExecution(t *testing.T) {
	s := schedule.New()
	var count int64
	_ = s.Add("halting", 10*time.Millisecond, func() error {
		atomic.AddInt64(&count, 1)
		return nil
	})
	s.Start()
	time.Sleep(35 * time.Millisecond)
	s.Stop()
	before := atomic.LoadInt64(&count)
	time.Sleep(30 * time.Millisecond)
	after := atomic.LoadInt64(&count)
	if after != before {
		t.Errorf("job continued after Stop: before=%d after=%d", before, after)
	}
}

func TestMultipleJobsRunConcurrently(t *testing.T) {
	s := schedule.New()
	var a, b int64
	_ = s.Add("a", 15*time.Millisecond, func() error { atomic.AddInt64(&a, 1); return nil })
	_ = s.Add("b", 15*time.Millisecond, func() error { atomic.AddInt64(&b, 1); return nil })
	s.Start()
	time.Sleep(50 * time.Millisecond)
	s.Stop()
	if atomic.LoadInt64(&a) == 0 || atomic.LoadInt64(&b) == 0 {
		t.Error("expected both jobs to have executed at least once")
	}
}
