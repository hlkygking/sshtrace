// Package export provides functionality to export SSH session data
// to various formats such as JSON and plain text.
package export

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"sshtrace/internal/session"
)

// Format represents an export output format.
type Format string

const (
	FormatJSON Format = "json"
	FormatText Format = "text"
)

// Exporter writes session data to an io.Writer in a given format.
type Exporter struct {
	format Format
	w      io.Writer
}

// New creates a new Exporter writing to w in the given format.
func New(w io.Writer, format Format) *Exporter {
	return &Exporter{format: format, w: w}
}

// Export serializes the session to the configured format.
func (e *Exporter) Export(s *session.Session) error {
	switch e.format {
	case FormatJSON:
		return e.exportJSON(s)
	case FormatText:
		return e.exportText(s)
	default:
		return fmt.Errorf("unsupported export format: %s", e.format)
	}
}

func (e *Exporter) exportJSON(s *session.Session) error {
	enc := json.NewEncoder(e.w)
	enc.SetIndent("", "  ")
	return enc.Encode(s)
}

func (e *Exporter) exportText(s *session.Session) error {
	_, err := fmt.Fprintf(e.w, "Session ID : %s\n", s.ID)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(e.w, "User       : %s\n", s.User)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(e.w, "Remote IP  : %s\n", s.RemoteIP)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(e.w, "Started    : %s\n", s.StartedAt.Format(time.RFC3339))
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(e.w, "Events     : %d\n", len(s.Events))
	if err != nil {
		return err
	}
	for i, ev := range s.Events {
		_, err = fmt.Fprintf(e.w, "  [%d] %s (%s) %q\n", i+1, ev.Timestamp.Format(time.RFC3339), ev.Type, ev.Data)
		if err != nil {
			return err
		}
	}
	return nil
}
