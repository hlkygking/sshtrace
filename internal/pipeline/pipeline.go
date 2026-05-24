// Package pipeline chains multiple session processors together,
// applying each in order to produce a final processed session.
package pipeline

import (
	"errors"
	"fmt"

	"github.com/sshtrace/sshtrace/internal/session"
)

// Processor is any type that can transform a session.
type Processor interface {
	Process(s *session.Session) (*session.Session, error)
}

// Pipeline applies a sequence of Processors to a session.
type Pipeline struct {
	steps []Processor
}

// New creates a Pipeline with the given processors applied in order.
func New(steps ...Processor) (*Pipeline, error) {
	if len(steps) == 0 {
		return nil, errors.New("pipeline: at least one processor is required")
	}
	return &Pipeline{steps: steps}, nil
}

// Run passes the session through each processor in sequence.
// Processing stops and the error is returned if any step fails.
func (p *Pipeline) Run(s *session.Session) (*session.Session, error) {
	if s == nil {
		return nil, errors.New("pipeline: session must not be nil")
	}
	current := s
	for i, step := range p.steps {
		result, err := step.Process(current)
		if err != nil {
			return nil, fmt.Errorf("pipeline: step %d failed: %w", i, err)
		}
		if result == nil {
			return nil, fmt.Errorf("pipeline: step %d returned nil session", i)
		}
		current = result
	}
	return current, nil
}

// Len returns the number of processors in the pipeline.
func (p *Pipeline) Len() int {
	return len(p.steps)
}
