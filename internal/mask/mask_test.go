package mask_test

import (
	"testing"
	"time"

	"sshtrace/internal/mask"
	"sshtrace/internal/session"
)

func makeSession(events []session.Event) session.Session {
	s := session.Session{
		ID:        "test-id",
		User:      "alice",
		StartedAt: time.Now(),
		Events:    events,
	}
	return s
}

func TestNewRequiresFields(t *testing.T) {
	_, err := mask.New(nil, "")
	if err == nil {
		t.Fatal("expected error for empty fields slice")
	}
}

func TestNewRejectsBlankField(t *testing.T) {
	_, err := mask.New([]string{"  "}, "")
	if err == nil {
		t.Fatal("expected error for blank field name")
	}
}

func TestDefaultPlaceholder(t *testing.T) {
	m, err := mask.New([]string{"password"}, "")
	if err != nil {
		t.Fatal(err)
	}
	s := makeSession([]session.Event{
		{Type: "password", Data: "s3cr3t"},
	})
	out := m.Apply(s)
	if out.Events[0].Data != "[MASKED]" {
		t.Errorf("expected [MASKED], got %q", out.Events[0].Data)
	}
}

func TestCustomPlaceholder(t *testing.T) {
	m, err := mask.New([]string{"token"}, "***")
	if err != nil {
		t.Fatal(err)
	}
	s := makeSession([]session.Event{
		{Type: "token", Data: "abc123"},
	})
	out := m.Apply(s)
	if out.Events[0].Data != "***" {
		t.Errorf("expected ***, got %q", out.Events[0].Data)
	}
}

func TestUnmaskedFieldsUntouched(t *testing.T) {
	m, err := mask.New([]string{"password"}, "")
	if err != nil {
		t.Fatal(err)
	}
	s := makeSession([]session.Event{
		{Type: "command", Data: "ls -la"},
		{Type: "password", Data: "hunter2"},
	})
	out := m.Apply(s)
	if out.Events[0].Data != "ls -la" {
		t.Errorf("command data should be unchanged, got %q", out.Events[0].Data)
	}
	if out.Events[1].Data != "[MASKED]" {
		t.Errorf("password data should be masked, got %q", out.Events[1].Data)
	}
}

func TestOriginalSessionUnmodified(t *testing.T) {
	m, err := mask.New([]string{"password"}, "")
	if err != nil {
		t.Fatal(err)
	}
	s := makeSession([]session.Event{
		{Type: "password", Data: "original"},
	})
	m.Apply(s)
	if s.Events[0].Data != "original" {
		t.Error("original session should not be modified")
	}
}

func TestMultipleFieldsMasked(t *testing.T) {
	m, err := mask.New([]string{"password", "token", "key"}, "")
	if err != nil {
		t.Fatal(err)
	}
	s := makeSession([]session.Event{
		{Type: "password", Data: "pw"},
		{Type: "token", Data: "tok"},
		{Type: "key", Data: "k"},
		{Type: "output", Data: "hello"},
	})
	out := m.Apply(s)
	for i := 0; i < 3; i++ {
		if out.Events[i].Data != "[MASKED]" {
			t.Errorf("event %d should be masked", i)
		}
	}
	if out.Events[3].Data != "hello" {
		t.Error("output event should not be masked")
	}
}
