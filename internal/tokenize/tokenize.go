// Package tokenize splits SSH session command events into discrete tokens
// (words/arguments) to enable fine-grained analysis and searching.
package tokenize

import (
	"errors"
	"strings"
	"unicode"

	"sshtrace/internal/session"
)

// Tokenizer splits event data into tokens.
type Tokenizer struct {
	// MinLength discards tokens shorter than this value. Zero means keep all.
	MinLength int
	// MaxTokens limits the number of tokens stored per event. Zero means unlimited.
	MaxTokens int
}

// New returns a Tokenizer with the given minimum token length and max tokens per event.
// MinLength must be >= 0 and MaxTokens must be >= 0.
func New(minLength, maxTokens int) (*Tokenizer, error) {
	if minLength < 0 {
		return nil, errors.New("tokenize: minLength must be >= 0")
	}
	if maxTokens < 0 {
		return nil, errors.New("tokenize: maxTokens must be >= 0")
	}
	return &Tokenizer{MinLength: minLength, MaxTokens: maxTokens}, nil
}

// Tokenize splits the given raw string into tokens, applying length and count limits.
func (t *Tokenizer) Tokenize(raw string) []string {
	fields := strings.FieldsFunc(raw, func(r rune) bool {
		return unicode.IsSpace(r) || r == '|' || r == ';' || r == '&'
	})

	var tokens []string
	for _, f := range fields {
		if t.MinLength > 0 && len(f) < t.MinLength {
			continue
		}
		tokens = append(tokens, f)
		if t.MaxTokens > 0 && len(tokens) >= t.MaxTokens {
			break
		}
	}
	return tokens
}

// Apply walks all events in s and stores tokenized data in each event's Meta map
// under the key "tokens" (space-joined). Returns an error if s is nil.
func (t *Tokenizer) Apply(s *session.Session) error {
	if s == nil {
		return errors.New("tokenize: session must not be nil")
	}
	for i := range s.Events {
		tokens := t.Tokenize(s.Events[i].Data)
		if s.Events[i].Meta == nil {
			s.Events[i].Meta = make(map[string]string)
		}
		s.Events[i].Meta["tokens"] = strings.Join(tokens, " ")
	}
	return nil
}
