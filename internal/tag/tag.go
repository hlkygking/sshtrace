// Package tag provides session tagging and label management for sshtrace.
// Tags allow operators to categorise sessions (e.g. "prod", "ci", "admin")
// and later filter or report on them.
package tag

import (
	"errors"
	"regexp"
	"sort"
	"strings"

	"github.com/sshtrace/sshtrace/internal/session"
)

var validTag = regexp.MustCompile(`^[a-z0-9_\-]{1,32}$`)

// ErrInvalidTag is returned when a tag name does not match the allowed pattern.
var ErrInvalidTag = errors.New("tag: invalid tag name; must match [a-z0-9_\\-]{1,32}")

// Tagger attaches and removes tags on sessions.
type Tagger struct{}

// New returns a new Tagger.
func New() *Tagger { return &Tagger{} }

// Validate reports whether name is a legal tag.
func Validate(name string) error {
	if !validTag.MatchString(name) {
		return ErrInvalidTag
	}
	return nil
}

// Add appends tags to s.Tags, deduplicating and sorting the result.
// Returns ErrInvalidTag if any supplied tag is malformed.
func (t *Tagger) Add(s *session.Session, tags ...string) error {
	for _, tag := range tags {
		if err := Validate(tag); err != nil {
			return err
		}
	}
	existing := make(map[string]struct{}, len(s.Tags))
	for _, tag := range s.Tags {
		existing[tag] = struct{}{}
	}
	for _, tag := range tags {
		existing[tag] = struct{}{}
	}
	result := make([]string, 0, len(existing))
	for tag := range existing {
		result = append(result, tag)
	}
	sort.Strings(result)
	s.Tags = result
	return nil
}

// Remove deletes the given tags from s.Tags.
func (t *Tagger) Remove(s *session.Session, tags ...string) {
	remove := make(map[string]struct{}, len(tags))
	for _, tag := range tags {
		remove[strings.ToLower(tag)] = struct{}{}
	}
	filtered := s.Tags[:0]
	for _, tag := range s.Tags {
		if _, skip := remove[tag]; !skip {
			filtered = append(filtered, tag)
		}
	}
	s.Tags = filtered
}

// Has reports whether s carries all of the supplied tags.
func Has(s *session.Session, tags ...string) bool {
	set := make(map[string]struct{}, len(s.Tags))
	for _, tag := range s.Tags {
		set[tag] = struct{}{}
	}
	for _, tag := range tags {
		if _, ok := set[tag]; !ok {
			return false
		}
	}
	return true
}
