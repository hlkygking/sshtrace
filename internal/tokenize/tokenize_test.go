package tokenize_test

import (
	"strings"
	"testing"
	"time"

	"sshtrace/internal/session"
	"sshtrace/internal/tokenize"
)

func makeSession(events ...string) *session.Session {
	s := session.New("user1", "127.0.0.1")
	for _, e := range events {
		s.AddEvent(session.Event{Kind: "output", Data: e, Timestamp: time.Now()})
	}
	return s
}

func TestNewInvalidMinLength(t *testing.T) {
	_, err := tokenize.New(-1, 0)
	if err == nil {
		t.Fatal("expected error for negative minLength")
	}
}

func TestNewInvalidMaxTokens(t *testing.T) {
	_, err := tokenize.New(0, -1)
	if err == nil {
		t.Fatal("expected error for negative maxTokens")
	}
}

func TestTokenizeSimple(t *testing.T) {
	tk, _ := tokenize.New(0, 0)
	tokens := tk.Tokenize("ls -la /home")
	if len(tokens) != 3 {
		t.Fatalf("expected 3 tokens, got %d", len(tokens))
	}
}

func TestTokenizePipeAndSemicolon(t *testing.T) {
	tk, _ := tokenize.New(0, 0)
	tokens := tk.Tokenize("cat file.txt | grep foo; echo done")
	if len(tokens) != 5 {
		t.Fatalf("expected 5 tokens, got %d: %v", len(tokens), tokens)
	}
}

func TestTokenizeMinLength(t *testing.T) {
	tk, _ := tokenize.New(3, 0)
	tokens := tk.Tokenize("ls -la /home")
	// "-la" is 3 chars (kept), "ls" is 2 (dropped)
	for _, tok := range tokens {
		if len(tok) < 3 {
			t.Errorf("token %q is shorter than minLength 3", tok)
		}
	}
}

func TestTokenizeMaxTokens(t *testing.T) {
	tk, _ := tokenize.New(0, 2)
	tokens := tk.Tokenize("a b c d e")
	if len(tokens) != 2 {
		t.Fatalf("expected 2 tokens due to maxTokens, got %d", len(tokens))
	}
}

func TestApplyNilSessionReturnsError(t *testing.T) {
	tk, _ := tokenize.New(0, 0)
	if err := tk.Apply(nil); err == nil {
		t.Fatal("expected error for nil session")
	}
}

func TestApplyStoresTokensInMeta(t *testing.T) {
	tk, _ := tokenize.New(0, 0)
	s := makeSession("sudo rm -rf /tmp", "echo hello")
	if err := tk.Apply(s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, ev := range s.Events {
		if ev.Meta == nil {
			t.Fatal("expected Meta to be set")
		}
		if _, ok := ev.Meta["tokens"]; !ok {
			t.Error("expected 'tokens' key in Meta")
		}
	}
}

func TestApplyTokensContent(t *testing.T) {
	tk, _ := tokenize.New(0, 0)
	s := makeSession("git commit -m 'fix bug'")
	_ = tk.Apply(s)
	tokens := s.Events[0].Meta["tokens"]
	if !strings.Contains(tokens, "git") {
		t.Errorf("expected 'git' in tokens, got %q", tokens)
	}
}
