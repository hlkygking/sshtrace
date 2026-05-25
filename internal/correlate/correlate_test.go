package correlate_test

import (
	"testing"
	"time"

	"sshtrace/internal/correlate"
	"sshtrace/internal/session"
)

func makeSession(user, ip string, start time.Time) *session.Session {
	s := session.New(user, ip)
	s.StartedAt = start
	return s
}

func TestDefaultOptions(t *testing.T) {
	opts := correlate.DefaultOptions()
	if opts.Strategy != correlate.ByUser {
		t.Fatalf("expected ByUser, got %v", opts.Strategy)
	}
	if opts.Window != 5*time.Minute {
		t.Fatalf("expected 5m window, got %v", opts.Window)
	}
}

func TestNewInvalidWindowReturnsError(t *testing.T) {
	_, err := correlate.New(correlate.Options{
		Strategy: correlate.ByTimeWindow,
		Window:   0,
	})
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestGroupByUserEmptyInput(t *testing.T) {
	c, _ := correlate.New(correlate.DefaultOptions())
	result := c.Group(nil)
	if len(result) != 0 {
		t.Fatalf("expected empty result, got %d clusters", len(result))
	}
}

func TestGroupByUser(t *testing.T) {
	now := time.Now()
	sessions := []*session.Session{
		makeSession("alice", "1.2.3.4", now),
		makeSession("bob", "1.2.3.5", now),
		makeSession("alice", "1.2.3.6", now),
	}
	c, _ := correlate.New(correlate.DefaultOptions())
	clusters := c.Group(sessions)
	if len(clusters) != 2 {
		t.Fatalf("expected 2 clusters, got %d", len(clusters))
	}
	// alice comes first (preserves insertion order)
	if len(clusters[0]) != 2 {
		t.Fatalf("expected alice cluster size 2, got %d", len(clusters[0]))
	}
	if len(clusters[1]) != 1 {
		t.Fatalf("expected bob cluster size 1, got %d", len(clusters[1]))
	}
}

func TestGroupByIP(t *testing.T) {
	now := time.Now()
	sessions := []*session.Session{
		makeSession("alice", "10.0.0.1", now),
		makeSession("bob", "10.0.0.2", now),
		makeSession("carol", "10.0.0.1", now),
	}
	c, _ := correlate.New(correlate.Options{Strategy: correlate.ByIP})
	clusters := c.Group(sessions)
	if len(clusters) != 2 {
		t.Fatalf("expected 2 clusters, got %d", len(clusters))
	}
}

func TestGroupByTimeWindowOverlapping(t *testing.T) {
	base := time.Now()
	sessions := []*session.Session{
		makeSession("alice", "1.1.1.1", base),
		makeSession("bob", "2.2.2.2", base.Add(2*time.Minute)),
		makeSession("carol", "3.3.3.3", base.Add(20*time.Minute)),
	}
	c, _ := correlate.New(correlate.Options{
		Strategy: correlate.ByTimeWindow,
		Window:   5 * time.Minute,
	})
	clusters := c.Group(sessions)
	// alice and bob overlap; carol does not
	if len(clusters) != 2 {
		t.Fatalf("expected 2 clusters, got %d", len(clusters))
	}
	if len(clusters[0]) != 2 {
		t.Fatalf("expected first cluster size 2, got %d", len(clusters[0]))
	}
}

func TestGroupByTimeWindowNoOverlap(t *testing.T) {
	base := time.Now()
	sessions := []*session.Session{
		makeSession("alice", "1.1.1.1", base),
		makeSession("bob", "2.2.2.2", base.Add(10*time.Minute)),
	}
	c, _ := correlate.New(correlate.Options{
		Strategy: correlate.ByTimeWindow,
		Window:   2 * time.Minute,
	})
	clusters := c.Group(sessions)
	if len(clusters) != 2 {
		t.Fatalf("expected 2 clusters, got %d", len(clusters))
	}
}
