package rotate_test

import (
	"testing"
	"time"

	"sshtrace/internal/rotate"
	"sshtrace/internal/session"
)

func makeSession(id string, age time.Duration) *session.Session {
	s := session.New(id, "testuser", "127.0.0.1")
	s.StartedAt = time.Now().Add(-age)
	return s
}

func TestDefaultPolicy(t *testing.T) {
	p := rotate.DefaultPolicy()
	if p.MaxAge <= 0 {
		t.Error("expected positive MaxAge")
	}
	if p.MaxCount <= 0 {
		t.Error("expected positive MaxCount")
	}
}

func TestNewInvalidMaxAge(t *testing.T) {
	_, err := rotate.New(rotate.Policy{MaxAge: -1})
	if err == nil {
		t.Error("expected error for negative MaxAge")
	}
}

func TestNewInvalidMaxCount(t *testing.T) {
	_, err := rotate.New(rotate.Policy{MaxCount: -1})
	if err == nil {
		t.Error("expected error for negative MaxCount")
	}
}

func TestRotateByAge(t *testing.T) {
	rot, _ := rotate.New(rotate.Policy{MaxAge: 24 * time.Hour})

	sessions := []*session.Session{
		makeSession("old1", 48*time.Hour),
		makeSession("old2", 36*time.Hour),
		makeSession("new1", 1*time.Hour),
	}

	keep, remove := rot.Apply(sessions)
	if len(keep) != 1 {
		t.Fatalf("expected 1 kept, got %d", len(keep))
	}
	if len(remove) != 2 {
		t.Fatalf("expected 2 removed, got %d", len(remove))
	}
	if keep[0].ID != "new1" {
		t.Errorf("expected new1 to be kept, got %s", keep[0].ID)
	}
}

func TestRotateByCount(t *testing.T) {
	rot, _ := rotate.New(rotate.Policy{MaxCount: 2})

	sessions := []*session.Session{
		makeSession("s1", 3*time.Hour),
		makeSession("s2", 2*time.Hour),
		makeSession("s3", 1*time.Hour),
	}

	keep, remove := rot.Apply(sessions)
	if len(keep) != 2 {
		t.Fatalf("expected 2 kept, got %d", len(keep))
	}
	if len(remove) != 1 {
		t.Fatalf("expected 1 removed, got %d", len(remove))
	}
	if remove[0].ID != "s1" {
		t.Errorf("expected s1 removed, got %s", remove[0].ID)
	}
}

func TestRotateNoLimits(t *testing.T) {
	rot, _ := rotate.New(rotate.Policy{})

	sessions := []*session.Session{
		makeSession("s1", 100*24*time.Hour),
		makeSession("s2", 1*time.Hour),
	}

	keep, remove := rot.Apply(sessions)
	if len(keep) != 2 {
		t.Errorf("expected all sessions kept, got %d", len(keep))
	}
	if len(remove) != 0 {
		t.Errorf("expected no sessions removed, got %d", len(remove))
	}
}

func TestRotateCombined(t *testing.T) {
	rot, _ := rotate.New(rotate.Policy{
		MaxAge:   10 * time.Hour,
		MaxCount: 1,
	})

	sessions := []*session.Session{
		makeSession("ancient", 48*time.Hour),
		makeSession("recent1", 2*time.Hour),
		makeSession("recent2", 1*time.Hour),
	}

	keep, remove := rot.Apply(sessions)
	if len(keep) != 1 {
		t.Fatalf("expected 1 kept, got %d", len(keep))
	}
	if keep[0].ID != "recent2" {
		t.Errorf("expected recent2 kept, got %s", keep[0].ID)
	}
	if len(remove) != 2 {
		t.Fatalf("expected 2 removed, got %d", len(remove))
	}
}
