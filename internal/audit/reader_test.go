package audit_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/sshtrace/sshtrace/internal/audit"
)

func writeRecords(t *testing.T, records []struct {
	kind   audit.EventKind
	user   string
	remote string
}) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	l := audit.New(&buf)
	for _, r := range records {
		s := makeSession()
		s.User = r.user
		s.RemoteIP = r.remote
		_ = l.LogSession(s, r.kind)
	}
	return &buf
}

func TestReadAll(t *testing.T) {
	buf := writeRecords(t, []struct {
		kind   audit.EventKind
		user   string
		remote string
	}{
		{audit.EventSessionStarted, "alice", "1.1.1.1"},
		{audit.EventSessionClosed, "bob", "2.2.2.2"},
	})

	recs, err := audit.Read(buf, audit.ReadOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(recs) != 2 {
		t.Errorf("expected 2 records, got %d", len(recs))
	}
}

func TestReadFilterByKind(t *testing.T) {
	buf := writeRecords(t, []struct {
		kind   audit.EventKind
		user   string
		remote string
	}{
		{audit.EventSessionStarted, "alice", "1.1.1.1"},
		{audit.EventCommandRun, "alice", "1.1.1.1"},
		{audit.EventSessionClosed, "alice", "1.1.1.1"},
	})

	recs, err := audit.Read(buf, audit.ReadOptions{Kind: audit.EventCommandRun})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(recs) != 1 {
		t.Errorf("expected 1 record, got %d", len(recs))
	}
}

func TestReadFilterByUser(t *testing.T) {
	buf := writeRecords(t, []struct {
		kind   audit.EventKind
		user   string
		remote string
	}{
		{audit.EventSessionStarted, "alice", "1.1.1.1"},
		{audit.EventSessionStarted, "bob", "2.2.2.2"},
	})

	recs, err := audit.Read(buf, audit.ReadOptions{User: "alice"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(recs) != 1 || recs[0].User != "alice" {
		t.Errorf("expected 1 alice record, got %v", recs)
	}
}

func TestReadFilterBySince(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)
	s := makeSession()
	_ = l.LogSession(s, audit.EventSessionStarted)

	future := time.Now().UTC().Add(time.Hour)
	recs, err := audit.Read(&buf, audit.ReadOptions{Since: future})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(recs) != 0 {
		t.Errorf("expected 0 records after future cutoff, got %d", len(recs))
	}
}

func TestReadInvalidJSON(t *testing.T) {
	buf := bytes.NewBufferString("not json\n")
	_, err := audit.Read(buf, audit.ReadOptions{})
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}
