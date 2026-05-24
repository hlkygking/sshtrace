// Package classify categorises SSH session events by command risk level.
package classify

import (
	"regexp"
	"strings"

	"sshtrace/internal/session"
)

// Level represents the risk classification of a session event.
type Level int

const (
	LevelSafe     Level = iota // routine, read-only commands
	LevelModerate              // potentially impactful commands
	LevelDangerous             // destructive or privilege-escalating commands
)

func (l Level) String() string {
	switch l {
	case LevelSafe:
		return "safe"
	case LevelModerate:
		return "moderate"
	case LevelDangerous:
		return "dangerous"
	default:
		return "unknown"
	}
}

var dangerousPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)\brm\s+-rf`),
	regexp.MustCompile(`(?i)\bsudo\b`),
	regexp.MustCompile(`(?i)\bchmod\s+777`),
	regexp.MustCompile(`(?i)\bdd\b.*of=`),
	regexp.MustCompile(`(?i)\b(shutdown|reboot|halt)\b`),
	regexp.MustCompile(`(?i)\bpasswd\b`),
}

var moderatePatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)\b(apt|yum|brew)\s+install`),
	regexp.MustCompile(`(?i)\bsystemctl\b`),
	regexp.MustCompile(`(?i)\bcurl\b.*\|.*sh`),
	regexp.MustCompile(`(?i)\bwget\b.*\|.*sh`),
	regexp.MustCompile(`(?i)\bchown\b`),
}

// Classify returns the risk Level for a single data string.
func Classify(data string) Level {
	for _, p := range dangerousPatterns {
		if p.MatchString(data) {
			return LevelDangerous
		}
	}
	for _, p := range moderatePatterns {
		if p.MatchString(data) {
			return LevelModerate
		}
	}
	return LevelSafe
}

// Result holds the classification outcome for a single event.
type Result struct {
	EventIndex int
	Data       string
	Level      Level
}

// Analyse classifies every event in the session and returns a slice of Results.
// Only output events (Kind == "output") are inspected; input events are skipped.
func Analyse(s *session.Session) []Result {
	if s == nil {
		return nil
	}
	var results []Result
	for i, e := range s.Events {
		if strings.ToLower(e.Kind) != "output" {
			continue
		}
		results = append(results, Result{
			EventIndex: i,
			Data:       e.Data,
			Level:      Classify(e.Data),
		})
	}
	return results
}
