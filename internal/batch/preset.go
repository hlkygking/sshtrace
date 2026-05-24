package batch

import (
	"strings"

	"sshtrace/internal/session"
)

// NoOpStep is a Step that does nothing and always succeeds.
// Useful as a placeholder in tests or pipelines under construction.
func NoOpStep(s *session.Session) error { return nil }

// UppercaseUserStep converts the session User field to upper-case.
// Intended primarily for demonstration and testing purposes.
func UppercaseUserStep(s *session.Session) error {
	s.User = strings.ToUpper(s.User)
	return nil
}

// PresetProcessor returns a ready-made Processor for common workflows.
// Recognised names:
//
//	"noop"  – single no-op step (useful for smoke tests)
//
Any unrecognised name falls back to the noop preset.
func PresetProcessor(name string) (*Processor, error) {
	switch strings.ToLower(name) {
	case "noop":
		return New(NoOpStep)
	default:
		return New(NoOpStep)
	}
}
