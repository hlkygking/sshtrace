// Package enrich augments session metadata with derived fields such as
// geo-location hints, hostname lookups, and session classification labels.
package enrich

import (
	"fmt"
	"net"
	"strings"

	"sshtrace/internal/session"
)

// Enricher adds derived metadata to a session.
type Enricher struct {
	resolveHostnames bool
	classifyCommands bool
}

// Option configures an Enricher.
type Option func(*Enricher)

// WithHostnameResolution enables reverse-DNS lookup for the session's remote IP.
func WithHostnameResolution() Option {
	return func(e *Enricher) { e.resolveHostnames = true }
}

// WithCommandClassification enables tagging events with a command category.
func WithCommandClassification() Option {
	return func(e *Enricher) { e.classifyCommands = true }
}

// New creates an Enricher with the supplied options.
func New(opts ...Option) *Enricher {
	e := &Enricher{}
	for _, o := range opts {
		o(e)
	}
	return e
}

// Apply enriches the session in-place and returns it.
func (e *Enricher) Apply(s *session.Session) (*session.Session, error) {
	if s == nil {
		return nil, fmt.Errorf("enrich: nil session")
	}

	if e.resolveHostnames && s.RemoteIP != "" {
		if host := resolveHost(s.RemoteIP); host != "" {
			s.Meta["hostname"] = host
		}
	}

	if e.classifyCommands {
		for i, ev := range s.Events {
			if ev.Type == "output" {
				continue
			}
			s.Events[i].Meta["category"] = classifyCommand(ev.Data)
		}
	}

	return s, nil
}

func resolveHost(ip string) string {
	hosts, err := net.LookupAddr(ip)
	if err != nil || len(hosts) == 0 {
		return ""
	}
	return strings.TrimSuffix(hosts[0], ".")
}

func classifyCommand(data string) string {
	d := strings.TrimSpace(strings.ToLower(data))
	switch {
	case strings.HasPrefix(d, "sudo") || strings.HasPrefix(d, "su "):
		return "privilege-escalation"
	case strings.HasPrefix(d, "ssh") || strings.HasPrefix(d, "scp") || strings.HasPrefix(d, "sftp"):
		return "remote-access"
	case strings.HasPrefix(d, "cat ") || strings.HasPrefix(d, "less ") || strings.HasPrefix(d, "more ") ||
		strings.HasPrefix(d, "tail ") || strings.HasPrefix(d, "head "):
		return "file-read"
	case strings.HasPrefix(d, "rm ") || strings.HasPrefix(d, "mv ") || strings.HasPrefix(d, "cp "):
		return "file-modify"
	case strings.HasPrefix(d, "curl") || strings.HasPrefix(d, "wget"):
		return "network-fetch"
	default:
		return "general"
	}
}
