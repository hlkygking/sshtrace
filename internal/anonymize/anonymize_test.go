package anonymize

import (
	"testing"
	"time"

	"github.com/sshtrace/sshtrace/internal/session"
)

func makeSession(user, ip string) *session.Session {
	s := session.New()
	s.User = user
	s.RemoteIP = ip
	s.StartedAt = time.Now()
	return s
}

func TestNewEmptySaltReturnsError(t *testing.T) {
	_, err := New("")
	if err == nil {
		t.Fatal("expected error for empty salt")
	}
}

func TestNewValidSalt(t *testing.T) {
	a, err := New("mysalt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a == nil {
		t.Fatal("expected non-nil Anonymizer")
	}
}

func TestApplyNilSessionReturnsError(t *testing.T) {
	a, _ := New("salt")
	_, err := a.Apply(nil)
	if err == nil {
		t.Fatal("expected error for nil session")
	}
}

func TestApplyReplacesUserAndIP(t *testing.T) {
	a, _ := New("salt")
	s := makeSession("alice", "192.168.1.10")
	out, err := a.Apply(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.User == "alice" {
		t.Error("expected user to be anonymised")
	}
	if out.RemoteIP == "192.168.1.10" {
		t.Error("expected IP to be anonymised")
	}
}

func TestApplyDoesNotModifyOriginal(t *testing.T) {
	a, _ := New("salt")
	s := makeSession("bob", "10.0.0.1")
	a.Apply(s) //nolint:errcheck
	if s.User != "bob" {
		t.Error("original session user should be unchanged")
	}
	if s.RemoteIP != "10.0.0.1" {
		t.Error("original session IP should be unchanged")
	}
}

func TestPseudonymIsStable(t *testing.T) {
	a, _ := New("stable-salt")
	s := makeSession("carol", "172.16.0.5")
	out1, _ := a.Apply(s)
	out2, _ := a.Apply(s)
	if out1.User != out2.User {
		t.Errorf("pseudonym not stable: %q vs %q", out1.User, out2.User)
	}
	if out1.RemoteIP != out2.RemoteIP {
		t.Errorf("IP pseudonym not stable: %q vs %q", out1.RemoteIP, out2.RemoteIP)
	}
}

func TestDifferentSaltsProduceDifferentPseudonyms(t *testing.T) {
	a1, _ := New("salt-one")
	a2, _ := New("salt-two")
	s := makeSession("dave", "10.1.2.3")
	out1, _ := a1.Apply(s)
	out2, _ := a2.Apply(s)
	if out1.User == out2.User {
		t.Error("different salts should produce different pseudonyms")
	}
}

func TestApplyPreservesStartedAt(t *testing.T) {
	a, _ := New("salt")
	now := time.Now()
	s := makeSession("eve", "10.0.0.2")
	s.StartedAt = now
	out, err := a.Apply(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !out.StartedAt.Equal(now) {
		t.Errorf("expected StartedAt to be preserved, got %v want %v", out.StartedAt, now)
	}
}
