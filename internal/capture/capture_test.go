package capture_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"sshtrace/internal/capture"
	"sshtrace/internal/session"
)

func newSession() *session.Session {
	return session.New("test-user", "127.0.0.1")
}

func TestWriteRecordsOutputEvent(t *testing.T) {
	sess := newSession()
	var buf bytes.Buffer
	c := capture.New(sess, strings.NewReader(""), &buf)

	payload := []byte("hello server")
	n, err := c.Write(payload)
	if err != nil {
		t.Fatalf("unexpected write error: %v", err)
	}
	if n != len(payload) {
		t.Fatalf("wrote %d bytes, want %d", n, len(payload))
	}
	if buf.String() != string(payload) {
		t.Errorf("underlying writer got %q, want %q", buf.String(), payload)
	}

	events := sess.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Direction != session.DirOutput {
		t.Errorf("expected DirOutput, got %v", events[0].Direction)
	}
	if string(events[0].Data) != string(payload) {
		t.Errorf("event data %q, want %q", events[0].Data, payload)
	}
}

func TestReadRecordsInputEvent(t *testing.T) {
	sess := newSession()
	payload := []byte("ls -la\n")
	c := capture.New(sess, bytes.NewReader(payload), &bytes.Buffer{})

	buf := make([]byte, 64)
	n, err := c.Read(buf)
	if err != nil {
		t.Fatalf("unexpected read error: %v", err)
	}
	if string(buf[:n]) != string(payload) {
		t.Errorf("read %q, want %q", buf[:n], payload)
	}

	events := sess.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Direction != session.DirInput {
		t.Errorf("expected DirInput, got %v", events[0].Direction)
	}
}

func TestWriteTimestampIsRecent(t *testing.T) {
	sess := newSession()
	before := time.Now()
	c := capture.New(sess, strings.NewReader(""), &bytes.Buffer{})
	_, _ = c.Write([]byte("ping"))
	after := time.Now()

	ts := sess.Events()[0].Timestamp
	if ts.Before(before) || ts.After(after) {
		t.Errorf("timestamp %v outside expected range [%v, %v]", ts, before, after)
	}
}
