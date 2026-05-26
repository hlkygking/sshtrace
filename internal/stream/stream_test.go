package stream_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/sshtrace/sshtrace/internal/session"
	"github.com/sshtrace/sshtrace/internal/stream"
)

func makeSession() *session.Session {
	s := session.New("alice", "10.0.0.1")
	return s
}

func makeEvent(kind, data string) session.Event {
	return session.Event{Kind: kind, Data: data, At: time.Now()}
}

func TestSubscribeAndPublish(t *testing.T) {
	b := stream.New()
	var got []string
	_ = b.Subscribe("s1", func(_ *session.Session, e session.Event) {
		got = append(got, e.Data)
	})
	sess := makeSession()
	b.Publish(sess, makeEvent("output", "hello"))
	b.Publish(sess, makeEvent("output", "world"))
	if len(got) != 2 || got[0] != "hello" || got[1] != "world" {
		t.Fatalf("expected [hello world], got %v", got)
	}
}

func TestSubscribeDuplicateReturnsError(t *testing.T) {
	b := stream.New()
	_ = b.Subscribe("s1", func(_ *session.Session, _ session.Event) {})
	err := b.Subscribe("s1", func(_ *session.Session, _ session.Event) {})
	if err == nil {
		t.Fatal("expected error for duplicate subscriber")
	}
}

func TestSubscribeEmptyNameReturnsError(t *testing.T) {
	b := stream.New()
	err := b.Subscribe("", func(_ *session.Session, _ session.Event) {})
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestSubscribeNilHandlerReturnsError(t *testing.T) {
	b := stream.New()
	err := b.Subscribe("s1", nil)
	if err == nil {
		t.Fatal("expected error for nil handler")
	}
}

func TestUnsubscribeRemovesHandler(t *testing.T) {
	b := stream.New()
	_ = b.Subscribe("s1", func(_ *session.Session, _ session.Event) {})
	_ = b.Unsubscribe("s1")
	if b.Count() != 0 {
		t.Fatal("expected 0 subscribers after unsubscribe")
	}
}

func TestUnsubscribeUnknownReturnsError(t *testing.T) {
	b := stream.New()
	err := b.Unsubscribe("nobody")
	if err == nil {
		t.Fatal("expected error for unknown subscriber")
	}
}

func TestCount(t *testing.T) {
	b := stream.New()
	if b.Count() != 0 {
		t.Fatal("expected 0 initially")
	}
	_ = b.Subscribe("a", func(_ *session.Session, _ session.Event) {})
	_ = b.Subscribe("b", func(_ *session.Session, _ session.Event) {})
	if b.Count() != 2 {
		t.Fatalf("expected 2, got %d", b.Count())
	}
}

func TestWriterHandler(t *testing.T) {
	var buf bytes.Buffer
	b := stream.New()
	_ = b.Subscribe("w", stream.WriterHandler(&buf))
	sess := makeSession()
	b.Publish(sess, makeEvent("output", "hi"))
	if !strings.Contains(buf.String(), "hi") {
		t.Fatalf("expected output to contain 'hi', got: %s", buf.String())
	}
}

func TestFilterHandlerAllowsMatchingKind(t *testing.T) {
	var count int
	h := stream.FilterHandler([]string{"output"}, func(_ *session.Session, _ session.Event) { count++ })
	sess := makeSession()
	h(sess, makeEvent("output", "x"))
	h(sess, makeEvent("input", "y"))
	if count != 1 {
		t.Fatalf("expected 1 call, got %d", count)
	}
}

func TestCountingHandler(t *testing.T) {
	var n int64
	b := stream.New()
	_ = b.Subscribe("c", stream.CountingHandler(&n, nil))
	sess := makeSession()
	b.Publish(sess, makeEvent("output", "a"))
	b.Publish(sess, makeEvent("output", "b"))
	if n != 2 {
		t.Fatalf("expected counter 2, got %d", n)
	}
}
