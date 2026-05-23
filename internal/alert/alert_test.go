package alert_test

import (
	"testing"
	"time"

	"sshtrace/internal/alert"
	"sshtrace/internal/session"
)

func makeSession(events []session.Event) *session.Session {
	s := session.New("user1", "192.168.1.1")
	for _, ev := range events {
		s.AddEvent(ev)
	}
	return s
}

func evt(data string) session.Event {
	return session.Event{Type: "output", Data: data, Timestamp: time.Now()}
}

func TestNoRulesNoAlerts(t *testing.T) {
	ev := alert.New(nil)
	s := makeSession([]session.Event{evt("ls -la")})
	if got := ev.Evaluate(s); len(got) != 0 {
		t.Fatalf("expected 0 alerts, got %d", len(got))
	}
}

func TestKeywordMatch(t *testing.T) {
	rules := []alert.Rule{
		{Name: "sudo", Level: alert.LevelWarn, Keyword: "sudo"},
	}
	ev := alert.New(rules)
	s := makeSession([]session.Event{evt("sudo rm file")})
	alerts := ev.Evaluate(s)
	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(alerts))
	}
	if alerts[0].Rule.Name != "sudo" {
		t.Errorf("unexpected rule name: %s", alerts[0].Rule.Name)
	}
}

func TestCaseInsensitiveMatch(t *testing.T) {
	rules := []alert.Rule{
		{Name: "rm-rf", Level: alert.LevelCrit, Keyword: "RM -RF"},
	}
	ev := alert.New(rules)
	s := makeSession([]session.Event{evt("rm -rf /")})
	if got := ev.Evaluate(s); len(got) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(got))
	}
}

func TestMultipleRulesMultipleAlerts(t *testing.T) {
	rules := []alert.Rule{
		{Name: "sudo", Level: alert.LevelWarn, Keyword: "sudo"},
		{Name: "passwd", Level: alert.LevelCrit, Keyword: "passwd"},
	}
	ev := alert.New(rules)
	s := makeSession([]session.Event{
		evt("sudo passwd root"),
	})
	if got := ev.Evaluate(s); len(got) != 2 {
		t.Fatalf("expected 2 alerts, got %d", len(got))
	}
}

func TestNoMatchReturnsEmpty(t *testing.T) {
	rules := []alert.Rule{
		{Name: "danger", Level: alert.LevelCrit, Keyword: "DROP TABLE"},
	}
	ev := alert.New(rules)
	s := makeSession([]session.Event{evt("ls"), evt("pwd")})
	if got := ev.Evaluate(s); len(got) != 0 {
		t.Fatalf("expected 0 alerts, got %d", len(got))
	}
}

func TestAlertContainsSessionMeta(t *testing.T) {
	rules := []alert.Rule{
		{Name: "wget", Level: alert.LevelInfo, Keyword: "wget"},
	}
	ev := alert.New(rules)
	s := makeSession([]session.Event{evt("wget http://example.com")})
	alerts := ev.Evaluate(s)
	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert")
	}
	if alerts[0].User != "user1" {
		t.Errorf("expected user user1, got %s", alerts[0].User)
	}
	if alerts[0].SessionID != s.ID {
		t.Errorf("session ID mismatch")
	}
}
