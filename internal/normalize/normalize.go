// Package normalize provides utilities for standardizing SSH session event
// data into a consistent format before storage or further processing.
package normalize

import (
	"strings"
	"unicode"

	"github.com/sshtrace/sshtrace/internal/session"
)

// Options controls which normalization steps are applied.
type Options struct {
	TrimWhitespace  bool
	LowercaseUser   bool
	NormalizeLineEndings bool
}

// DefaultOptions returns a sensible default normalization configuration.
func DefaultOptions() Options {
	return Options{
		TrimWhitespace:       true,
		LowercaseUser:        false,
		NormalizeLineEndings: true,
	}
}

// Normalizer applies normalization rules to sessions.
type Normalizer struct {
	opts Options
}

// New creates a Normalizer with the given options.
func New(opts Options) *Normalizer {
	return &Normalizer{opts: opts}
}

// Apply normalizes the session in place and returns it.
// Returns an error if the session is nil.
func (n *Normalizer) Apply(s *session.Session) (*session.Session, error) {
	if s == nil {
		return nil, fmt.Errorf("normalize: session must not be nil")
	}

	if n.opts.LowercaseUser {
		s.User = strings.ToLower(s.User)
	}

	for i := range s.Events {
		data := s.Events[i].Data

		if n.opts.NormalizeLineEndings {
			data = strings.ReplaceAll(data, "\r\n", "\n")
			data = strings.ReplaceAll(data, "\r", "\n")
		}

		if n.opts.TrimWhitespace {
			data = strings.TrimRightFunc(data, unicode.IsSpace)
		}

		s.Events[i].Data = data
	}

	return s, nil
}
