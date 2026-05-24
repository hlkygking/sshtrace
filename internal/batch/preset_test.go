package batch_test

import (
	"strings"
	"testing"

	"sshtrace/internal/batch"
)

func TestPresetNoopReturnsProcessor(t *testing.T) {
	p, err := batch.PresetProcessor("noop")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil processor")
	}
}

func TestPresetUnknownFallsBackToNoop(t *testing.T) {
	p, err := batch.PresetProcessor("unknown-preset")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil processor")
	}
}

func TestUppercaseUserStep(t *testing.T) {
	s := makeSession("alice", "1.2.3.4")
	if err := batch.UppercaseUserStep(s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.User != strings.ToUpper("alice") {
		t.Errorf("expected upper-case user, got %q", s.User)
	}
}
