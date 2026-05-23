// Package alert provides rule-based alerting for SSH session events.
package alert

import (
	"strings"
	"time"

	"sshtrace/internal/session"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelCrit  Level = "crit"
)

// Rule defines a condition that triggers an alert.
type Rule struct {
	Name    string
	Level   Level
	Keyword string // case-insensitive substring match against event data
}

// Alert is produced when a Rule matches an event.
type Alert struct {
	Rule      Rule
	SessionID string
	User      string
	Triggered time.Time
	MatchedData string
}

// Evaluator checks sessions against a set of rules.
type Evaluator struct {
	rules []Rule
}

// New creates an Evaluator with the provided rules.
func New(rules []Rule) *Evaluator {
	return &Evaluator{rules: rules}
}

// Evaluate scans all events in a session and returns any triggered alerts.
func (e *Evaluator) Evaluate(s *session.Session) []Alert {
	var alerts []Alert
	for _, ev := range s.Events {
		for _, rule := range e.rules {
			if strings.Contains(
				strings.ToLower(ev.Data),
				strings.ToLower(rule.Keyword),
			) {
				alerts = append(alerts, Alert{
					Rule:        rule,
					SessionID:   s.ID,
					User:        s.User,
					Triggered:   ev.Timestamp,
					MatchedData: ev.Data,
				})
			}
		}
	}
	return alerts
}
