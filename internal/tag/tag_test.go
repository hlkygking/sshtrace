package tag_test

import (
	"testing"
	"time"

	"github.com/sshtrace/sshtrace/internal/session"
	"github.com/sshtrace/sshtrace/internal/tag"
)

func makeSession() *session.Session {
	s, _ := session.New("user1", "192.168.1.1")
	s.StartedAt = time.Now()
	return s
}

func TestValidate(t *testing.T) {
	valid := []string{"prod", "ci-runner", "admin_ops", "a", strings.Repeat("x", 32)}
	for _, v := range valid {
		if err := tag.Validate(v); err != nil {
			t.Errorf("expected %q to be valid, got %v", v, err)
		}
	}
	invalid := []string{"", "UPPER", "has space", strings.Repeat("x", 33), "special!"}
	for _, v := range invalid {
		if err := tag.Validate(v); err == nil {
			t.Errorf("expected %q to be invalid", v)
		}
	}
}

func TestAddDeduplicates(t *testing.T) {
	s := makeSession()
	tgr := tag.New()
	_ = tgr.Add(s, "prod", "ci")
	_ = tgr.Add(s, "prod", "staging")
	if len(s.Tags) != 3 {
		t.Fatalf("expected 3 tags, got %d: %v", len(s.Tags), s.Tags)
	}
}

func TestAddSorts(t *testing.T) {
	s := makeSession()
	tgr := tag.New()
	_ = tgr.Add(s, "zzz", "aaa", "mmm")
	if s.Tags[0] != "aaa" || s.Tags[2] != "zzz" {
		t.Fatalf("tags not sorted: %v", s.Tags)
	}
}

func TestAddInvalidReturnsError(t *testing.T) {
	s := makeSession()
	tgr := tag.New()
	if err := tgr.Add(s, "INVALID"); err == nil {
		t.Fatal("expected error for invalid tag")
	}
	if len(s.Tags) != 0 {
		t.Fatal("tags should not be modified on error")
	}
}

func TestRemove(t *testing.T) {
	s := makeSession()
	tgr := tag.New()
	_ = tgr.Add(s, "prod", "ci", "staging")
	tgr.Remove(s, "ci")
	if tag.Has(s, "ci") {
		t.Fatal("ci should have been removed")
	}
	if !tag.Has(s, "prod", "staging") {
		t.Fatal("prod and staging should remain")
	}
}

func TestHasAll(t *testing.T) {
	s := makeSession()
	tgr := tag.New()
	_ = tgr.Add(s, "prod", "ci")
	if !tag.Has(s, "prod", "ci") {
		t.Fatal("expected Has to return true")
	}
	if tag.Has(s, "prod", "missing") {
		t.Fatal("expected Has to return false when one tag is absent")
	}
}

func TestHasEmpty(t *testing.T) {
	s := makeSession()
	if !tag.Has(s) {
		t.Fatal("Has with no args should always return true")
	}
}
