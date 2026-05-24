// Package batch provides utilities for processing multiple SSH sessions
// in bulk, applying a transformation or predicate to each one.
package batch

import (
	"errors"
	"fmt"

	"sshtrace/internal/session"
)

// Processor applies an operation to a slice of sessions.
type Processor struct {
	steps []Step
}

// Step is a function that transforms a session, returning an error on failure.
type Step func(s *session.Session) error

// Result holds the outcome of processing a single session.
type Result struct {
	SessionID string
	Err       error
}

// New creates a Processor with the provided steps.
// At least one step must be supplied.
func New(steps ...Step) (*Processor, error) {
	if len(steps) == 0 {
		return nil, errors.New("batch: at least one step is required")
	}
	return &Processor{steps: steps}, nil
}

// Run applies all steps to each session in order.
// Processing continues even if individual sessions fail.
// Results are returned in the same order as the input sessions.
func (p *Processor) Run(sessions []*session.Session) []Result {
	results := make([]Result, len(sessions))
	for i, s := range sessions {
		if s == nil {
			results[i] = Result{SessionID: "", Err: fmt.Errorf("batch: session at index %d is nil", i)}
			continue
		}
		var err error
		for _, step := range p.steps {
			if err = step(s); err != nil {
				break
			}
		}
		results[i] = Result{SessionID: s.ID, Err: err}
	}
	return results
}

// Errors returns only the results that contain an error.
func Errors(results []Result) []Result {
	var errs []Result
	for _, r := range results {
		if r.Err != nil {
			errs = append(errs, r)
		}
	}
	return errs
}
