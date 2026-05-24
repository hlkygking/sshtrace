package batch_test

import (
	"errors"
	"testing"
	"time"

	"sshtrace/internal/batch"
	"sshtrace/internal/session"
)

func makeSession(user, ip string) *session.Session {
	s := session.New(user, ip)
	s.StartedAt = time.Now()
	return s
}

func TestNewRequiresStep(t *testing.T) {
	_, err := batch.New()
	if err == nil {
		t.Fatal("expected error when no steps provided")
	}
}

func TestRunNoOpReturnsNoErrors(t *testing.T) {
	p, err := batch.New(batch.NoOpStep)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sessions := []*session.Session{makeSession("alice", "1.2.3.4"), makeSession("bob", "5.6.7.8")}
	results := p.Run(sessions)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if r.Err != nil {
			t.Errorf("unexpected error for session %s: %v", r.SessionID, r.Err)
		}
	}
}

func TestRunNilSessionReturnsError(t *testing.T) {
	p, _ := batch.New(batch.NoOpStep)
	results := p.Run([]*session.Session{nil})
	if results[0].Err == nil {
		t.Fatal("expected error for nil session")
	}
}

func TestRunStepErrorStopsChain(t *testing.T) {
	var secondCalled bool
	failStep := func(s *session.Session) error { return errors.New("step failed") }
	secondStep := func(s *session.Session) error { secondCalled = true; return nil }

	p, _ := batch.New(failStep, secondStep)
	p.Run([]*session.Session{makeSession("carol", "9.9.9.9")})
	if secondCalled {
		t.Fatal("second step should not have been called after first step error")
	}
}

func TestRunAppliesStepsInOrder(t *testing.T) {
	var order []int
	step := func(n int) batch.Step {
		return func(s *session.Session) error { order = append(order, n); return nil }
	}
	p, _ := batch.New(step(1), step(2), step(3))
	p.Run([]*session.Session{makeSession("dave", "0.0.0.0")})
	for i, v := range order {
		if v != i+1 {
			t.Fatalf("wrong order at index %d: got %d", i, v)
		}
	}
}

func TestErrorsFiltersResults(t *testing.T) {
	failStep := func(s *session.Session) error {
		if s.User == "bad" {
			return errors.New("bad user")
		}
		return nil
	}
	p, _ := batch.New(failStep)
	sessions := []*session.Session{makeSession("good", "1.1.1.1"), makeSession("bad", "2.2.2.2")}
	results := p.Run(sessions)
	errs := batch.Errors(results)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error result, got %d", len(errs))
	}
	if errs[0].SessionID != sessions[1].ID {
		t.Errorf("wrong session ID in error result")
	}
}
