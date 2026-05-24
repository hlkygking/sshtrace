package notify

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func baseMsg() Message {
	return Message{
		SessionID: "sess-1",
		User:      "alice",
		Event:     "alert",
		Detail:    "suspicious command",
		Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
	}
}

func TestLogNotifierWritesOutput(t *testing.T) {
	var buf bytes.Buffer
	n := NewLog(&buf)
	if err := n.Send(baseMsg()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "alice") {
		t.Errorf("expected user in output, got: %s", out)
	}
	if !strings.Contains(out, "suspicious command") {
		t.Errorf("expected detail in output, got: %s", out)
	}
}

func TestLogNotifierIncludesTimestamp(t *testing.T) {
	var buf bytes.Buffer
	n := NewLog(&buf)
	_ = n.Send(baseMsg())
	if !strings.Contains(buf.String(), "2024-01-15") {
		t.Errorf("expected timestamp in output, got: %s", buf.String())
	}
}

func TestWebhookNotifierSendsJSON(t *testing.T) {
	var received Message
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			http.Error(w, "bad json", 400)
			return
		}
		w.WriteHeader(200)
	}))
	defer ts.Close()

	n := NewWebhook(ts.URL)
	msg := baseMsg()
	if err := n.Send(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.User != msg.User {
		t.Errorf("expected user %q, got %q", msg.User, received.User)
	}
	if received.SessionID != msg.SessionID {
		t.Errorf("expected session_id %q, got %q", msg.SessionID, received.SessionID)
	}
}

func TestWebhookNotifierErrorOnBadStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer ts.Close()

	n := NewWebhook(ts.URL)
	if err := n.Send(baseMsg()); err == nil {
		t.Error("expected error on 500 response")
	}
}

func TestWebhookNotifierErrorOnBadURL(t *testing.T) {
	n := NewWebhook("http://127.0.0.1:1")
	if err := n.Send(baseMsg()); err == nil {
		t.Error("expected error on unreachable URL")
	}
}
