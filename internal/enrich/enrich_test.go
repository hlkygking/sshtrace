package enrich_test

import (
	"testing"
	"time"

	"sshtrace/internal/enrich"
	"sshtrace/internal/session"
)

func makeSession(ip string) *session.Session {
	s := session.New("alice", ip)
	return s
}

func addEvent(s *session.Session, typ, data string) {
	s.AddEvent(session.Event{
		Type:      typ,
		Data:      data,
		Timestamp: time.Now(),
		Meta:      map[string]string{},
	})
}

func TestApplyNilSessionReturnsError(t *testing.T) {
	e := enrich.New()
	_, err := e.Apply(nil)
	if err == nil {
		t.Fatal("expected error for nil session")
	}
}

func TestApplyNoOptionsIsNoop(t *testing.T) {
	s := makeSession("10.0.0.1")
	e := enrich.New()
	out, err := e.Apply(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != s {
		t.Fatal("expected same session pointer")
	}
}

func TestCommandClassificationPrivilegeEscalation(t *testing.T) {
	s := makeSession("10.0.0.1")
	addEvent(s, "input", "sudo rm -rf /tmp/foo")

	e := enrich.New(enrich.WithCommandClassification())
	out, err := e.Apply(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := out.Events[0].Meta["category"]
	if got != "privilege-escalation" {
		t.Errorf("expected privilege-escalation, got %q", got)
	}
}

func TestCommandClassificationNetworkFetch(t *testing.T) {
	s := makeSession("10.0.0.1")
	addEvent(s, "input", "curl https://example.com")

	e := enrich.New(enrich.WithCommandClassification())
	out, _ := e.Apply(s)
	if got := out.Events[0].Meta["category"]; got != "network-fetch" {
		t.Errorf("expected network-fetch, got %q", got)
	}
}

func TestCommandClassificationGeneral(t *testing.T) {
	s := makeSession("10.0.0.1")
	addEvent(s, "input", "ls -la")

	e := enrich.New(enrich.WithCommandClassification())
	out, _ := e.Apply(s)
	if got := out.Events[0].Meta["category"]; got != "general" {
		t.Errorf("expected general, got %q", got)
	}
}

func TestOutputEventsAreNotClassified(t *testing.T) {
	s := makeSession("10.0.0.1")
	addEvent(s, "output", "total 0")

	e := enrich.New(enrich.WithCommandClassification())
	out, _ := e.Apply(s)
	if _, ok := out.Events[0].Meta["category"]; ok {
		t.Error("output events should not receive a category")
	}
}

func TestHostnameResolutionEmptyIPIsSkipped(t *testing.T) {
	s := makeSession("")
	e := enrich.New(enrich.WithHostnameResolution())
	_, err := e.Apply(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := s.Meta["hostname"]; ok {
		t.Error("should not set hostname for empty IP")
	}
}
