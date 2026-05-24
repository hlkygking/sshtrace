// Package mask provides field-level masking for SSH session events,
// replacing sensitive values with a configurable placeholder string.
package mask

import (
	"fmt"
	"strings"

	"sshtrace/internal/session"
)

const defaultPlaceholder = "[MASKED]"

// Masker replaces the content of nominated event fields with a placeholder.
type Masker struct {
	fields      map[string]bool
	placeholder string
}

// New returns a Masker that will mask the given field names.
// placeholder is the string written in place of the original value;
// pass an empty string to use the default "[MASKED]".
func New(fields []string, placeholder string) (*Masker, error) {
	if len(fields) == 0 {
		return nil, fmt.Errorf("mask: at least one field name is required")
	}
	ph := placeholder
	if ph == "" {
		ph = defaultPlaceholder
	}
	set := make(map[string]bool, len(fields))
	for _, f := range fields {
		f = strings.TrimSpace(f)
		if f == "" {
			return nil, fmt.Errorf("mask: field name must not be blank")
		}
		set[f] = true
	}
	return &Masker{fields: set, placeholder: ph}, nil
}

// Apply returns a shallow copy of s with matching event Data fields replaced
// by the placeholder. The original session is not modified.
func (m *Masker) Apply(s session.Session) session.Session {
	out := s
	out.Events = make([]session.Event, len(s.Events))
	for i, ev := range s.Events {
		if m.fields[ev.Type] {
			ev.Data = m.placeholder
		}
		out.Events[i] = ev
	}
	return out
}
