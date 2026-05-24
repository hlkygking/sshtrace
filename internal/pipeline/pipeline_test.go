package pipeline_test

import (
	"errors"
	"testing"
	"time"

	"github.com/sshtrace/sshtrace/internal/pipeline"
	"github.com/sshtrace/sshtrace/internal/session"
)

func makeSession() *session.Session {
	s := session.New("alice", "192.168.1.1")
	s.AddEvent(session.Event{Kind: "output", Data: "hello", Timestamp: time.Now()})
	return s
}

func TestNewRequiresAtLeastOneStep(t *testing.T) {
	_, err := pipeline.New()
	if err == nil {
		t.Fatal("expected error for empty pipeline")
	}
}

func TestRunNilSessionReturnsError(t *testing.T) {
	p, _ := pipeline.New(pipeline.NoOp())
	_, err := p.Run(nil)
	if err == nil {
		t.Fatal("expected error for nil session")
	}
}

func TestRunPassesThroughNoOp(t *testing.T) {
	s := makeSession()
	p, _ := pipeline.New(pipeline.NoOp())
	out, err := p.Run(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.ID != s.ID {
		t.Errorf("expected same session ID, got %q", out.ID)
	}
}

func TestRunAppliesStepsInOrder(t *testing.T) {
	var order []int
	step := func(n int) pipeline.Processor {
		return pipeline.ProcessorFunc(func(s *session.Session) (*session.Session, error) {
			order = append(order, n)
			return s, nil
		})
	}
	p, _ := pipeline.New(step(1), step(2), step(3))
	p.Run(makeSession())
	for i, v := range order {
		if v != i+1 {
			t.Errorf("step %d ran at position %d", v, i)
		}
	}
}

func TestRunStopsOnError(t *testing.T) {
	ran := false
	failing := pipeline.ProcessorFunc(func(s *session.Session) (*session.Session, error) {
		return nil, errors.New("boom")
	})
	after := pipeline.ProcessorFunc(func(s *session.Session) (*session.Session, error) {
		ran = true
		return s, nil
	})
	p, _ := pipeline.New(failing, after)
	_, err := p.Run(makeSession())
	if err == nil {
		t.Fatal("expected error from failing step")
	}
	if ran {
		t.Error("step after failure should not have run")
	}
}

func TestLenReturnsStepCount(t *testing.T) {
	p, _ := pipeline.New(pipeline.NoOp(), pipeline.NoOp(), pipeline.NoOp())
	if p.Len() != 3 {
		t.Errorf("expected 3, got %d", p.Len())
	}
}
