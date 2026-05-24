package classify_test

import (
	"testing"
	"time"

	"sshtrace/internal/classify"
	"sshtrace/internal/session"
)

func makeSession(events []session.Event) *session.Session {
	s := &session.Session{
		ID:        "test-id",
		User:      "alice",
		IP:        "127.0.0.1",
		StartTime: time.Now(),
		Events:    events,
	}
	return s
}

func TestClassifySafe(t *testing.T) {
	for _, cmd := range []string{"ls -la", "cat /etc/hosts", "echo hello", "pwd"} {
		if got := classify.Classify(cmd); got != classify.LevelSafe {
			t.Errorf("Classify(%q) = %v, want LevelSafe", cmd, got)
		}
	}
}

func TestClassifyModerate(t *testing.T) {
	for _, cmd := range []string{
		"apt install nginx",
		"systemctl restart ssh",
		"curl https://example.com/setup.sh | sh",
		"chown root:root /etc/config",
	} {
		if got := classify.Classify(cmd); got != classify.LevelModerate {
			t.Errorf("Classify(%q) = %v, want LevelModerate", cmd, got)
		}
	}
}

func TestClassifyDangerous(t *testing.T) {
	for _, cmd := range []string{
		"sudo rm -rf /",
		"chmod 777 /etc/shadow",
		"dd if=/dev/zero of=/dev/sda",
		"shutdown -h now",
		"passwd root",
	} {
		if got := classify.Classify(cmd); got != classify.LevelDangerous {
			t.Errorf("Classify(%q) = %v, want LevelDangerous", cmd, got)
		}
	}
}

func TestLevelString(t *testing.T) {
	cases := map[classify.Level]string{
		classify.LevelSafe:      "safe",
		classify.LevelModerate:  "moderate",
		classify.LevelDangerous: "dangerous",
	}
	for lvl, want := range cases {
		if got := lvl.String(); got != want {
			t.Errorf("Level(%d).String() = %q, want %q", lvl, got, want)
		}
	}
}

func TestAnalyseNilSession(t *testing.T) {
	if results := classify.Analyse(nil); results != nil {
		t.Errorf("expected nil, got %v", results)
	}
}

func TestAnalyseSkipsInputEvents(t *testing.T) {
	s := makeSession([]session.Event{
		{Kind: "input", Data: "sudo rm -rf /"},
	})
	if results := classify.Analyse(s); len(results) != 0 {
		t.Errorf("expected 0 results for input-only session, got %d", len(results))
	}
}

func TestAnalyseReturnsCorrectLevels(t *testing.T) {
	s := makeSession([]session.Event{
		{Kind: "output", Data: "ls -la"},
		{Kind: "output", Data: "apt install curl"},
		{Kind: "output", Data: "sudo passwd alice"},
	})
	results := classify.Analyse(s)
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	expected := []classify.Level{classify.LevelSafe, classify.LevelModerate, classify.LevelDangerous}
	for i, r := range results {
		if r.Level != expected[i] {
			t.Errorf("result[%d]: got %v, want %v", i, r.Level, expected[i])
		}
	}
}
