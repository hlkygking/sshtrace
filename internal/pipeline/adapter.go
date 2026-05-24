package pipeline

import (
	"github.com/sshtrace/sshtrace/internal/session"
)

// ProcessorFunc is a function adapter that implements Processor.
type ProcessorFunc func(s *session.Session) (*session.Session, error)

// Process implements Processor by calling the underlying function.
func (f ProcessorFunc) Process(s *session.Session) (*session.Session, error) {
	return f(s)
}

// NoOp returns a Processor that returns the session unchanged.
// Useful as a placeholder or in tests.
func NoOp() Processor {
	return ProcessorFunc(func(s *session.Session) (*session.Session, error) {
		return s, nil
	})
}
