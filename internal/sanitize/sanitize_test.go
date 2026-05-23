package sanitize_test

import (
	"testing"
	"time"

	"sshtrace/internal/sanitize"
	"sshtrace/internal/session"
)

func mustNew(t *testing.T, extra ...string) *sanitize.Redactor {
	t.Helper()
	r, err := sanitize.New(extra...)
	if err != nil {
		t.Fatalf("sanitize.New: %v", err)
	}
	return r
}

func TestScrubPassword(t *testing.T) {
	r := mustNew(t)
	out := r.Scrub("mysql -u root --password=s3cr3t")
	if contains(out, "s3cr3t") {
		t.Errorf("expected password to be redacted, got: %s", out)
	}
}

func TestScrubToken(t *testing.T) {
	r := mustNew(t)
	out := r.Scrub("curl -H 'token: ghp_abc123XYZ'")
	if contains(out, "ghp_abc123XYZ") {
		t.Errorf("expected token to be redacted, got: %s", out)
	}
}

func TestScrubNoSensitiveData(t *testing.T) {
	r := mustNew(t)
	input := "ls -la /var/log"
	out := r.Scrub(input)
	if out != input {
		t.Errorf("expected unchanged output, got: %s", out)
	}
}

func TestScrubCustomPattern(t *testing.T) {
	r := mustNew(t, `(?i)ssn[=:\s]+\d{3}-\d{2}-\d{4}`)
	out := r.Scrub("ssn: 123-45-6789")
	if contains(out, "123-45-6789") {
		t.Errorf("expected SSN to be redacted, got: %s", out)
	}
}

func TestInvalidPatternReturnsError(t *testing.T) {
	_, err := sanitize.New(`[invalid`)
	if err == nil {
		t.Error("expected error for invalid regex, got nil")
	}
}

func TestScrubSession(t *testing.T) {
	r := mustNew(t)
	s := session.Session{
		ID:   "test-session",
		User: "alice",
		Events: []session.Event{
			{Type: session.EventOutput, Data: "echo hello", Timestamp: time.Now()},
			{Type: session.EventInput, Data: "password=hunter2", Timestamp: time.Now()},
		},
	}
	clean := r.ScrubSession(s)
	if contains(clean.Events[1].Data, "hunter2") {
		t.Errorf("expected password redacted in session event, got: %s", clean.Events[1].Data)
	}
	// original must be untouched
	if !contains(s.Events[1].Data, "hunter2") {
		t.Error("original session should not be modified")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		len(s) > 0 && containsStr(s, substr))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
