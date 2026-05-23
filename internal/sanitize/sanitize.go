// Package sanitize provides utilities for redacting sensitive data
// from SSH session events before storage or export.
package sanitize

import (
	"regexp"
	"strings"

	"sshtrace/internal/session"
)

// Redactor holds compiled patterns used to scrub sensitive content.
type Redactor struct {
	patterns []*regexp.Regexp
	replacement string
}

// defaultPatterns are common sensitive data patterns.
var defaultPatterns = []string{
	`(?i)password[=:\s]+\S+`,
	`(?i)passwd[=:\s]+\S+`,
	`(?i)secret[=:\s]+\S+`,
	`(?i)token[=:\s]+\S+`,
	`(?i)api[_-]?key[=:\s]+\S+`,
	`-----BEGIN [A-Z ]+-----[\s\S]+?-----END [A-Z ]+-----`,
}

// New creates a Redactor with the default sensitive-data patterns.
// Additional custom patterns (as regular expressions) may be supplied.
func New(extra ...string) (*Redactor, error) {
	allPatterns := append(defaultPatterns, extra...)
	compiled := make([]*regexp.Regexp, 0, len(allPatterns))
	for _, p := range allPatterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, err
		}
		compiled = append(compiled, re)
	}
	return &Redactor{
		patterns:    compiled,
		replacement: "[REDACTED]",
	}, nil
}

// Scrub returns a copy of text with all sensitive patterns replaced.
func (r *Redactor) Scrub(text string) string {
	for _, re := range r.patterns {
		text = re.ReplaceAllString(text, r.replacement)
	}
	return text
}

// ScrubSession returns a deep copy of s with all event data scrubbed.
func (r *Redactor) ScrubSession(s session.Session) session.Session {
	clean := s
	clean.Events = make([]session.Event, len(s.Events))
	for i, ev := range s.Events {
		ev.Data = r.Scrub(strings.TrimRight(ev.Data, "\n"))
		clean.Events[i] = ev
	}
	return clean
}
