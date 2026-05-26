package stream_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/sshtrace/sshtrace/internal/stream"
)

func TestFilterHandlerEmptyKindsForwardsAll(t *testing.T) {
	var count int
	h := stream.FilterHandler(nil, func(_ interface{ }, _ interface{ }) { count++ })
	// Use the real types via the broker.
	b := stream.New()
	_ = b.Subscribe("f", stream.FilterHandler(nil, func(s interface{}, e interface{}) {
		count++
	}))
	// Actually test with real session.Event types:
	var count2 int
	sess := makeSession()
	h2 := stream.FilterHandler([]string{}, func(_ interface{}, _ interface{}) { count2++ })
	_ = h2 // suppress unused warning — test via broker below

	b2 := stream.New()
	_ = b2.Subscribe("all", stream.FilterHandler([]string{}, func(s interface{}, e interface{}) {
		count2++
	}))
	b2.Publish(sess, makeEvent("output", "a"))
	b2.Publish(sess, makeEvent("input", "b"))
	if count2 != 2 {
		t.Fatalf("expected 2 events forwarded, got %d", count2)
	}
}

func TestWriterHandlerIncludesSessionID(t *testing.T) {
	var buf bytes.Buffer
	sess := makeSession()
	h := stream.WriterHandler(&buf)
	h(sess, makeEvent("output", "data"))
	if !strings.Contains(buf.String(), sess.ID) {
		t.Fatalf("expected session ID %q in output: %s", sess.ID, buf.String())
	}
}

func TestCountingHandlerWithInnerHandler(t *testing.T) {
	var n int64
	var inner int
	h := stream.CountingHandler(&n, func(_ interface{}, _ interface{}) { inner++ })
	sess := makeSession()
	h(sess, makeEvent("output", "x"))
	if n != 1 {
		t.Fatalf("expected counter 1, got %d", n)
	}
	if inner != 1 {
		t.Fatalf("expected inner called once, got %d", inner)
	}
}
