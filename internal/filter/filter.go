package filter

import (
	"strings"
	"time"

	"sshtrace/internal/session"
)

// Criteria holds the parameters used to filter sessions.
type Criteria struct {
	UserName  string
	RemoteIP  string
	Since     time.Time
	Until     time.Time
	MinEvents int
}

// Filter returns the subset of sessions that match all non-zero criteria.
func Filter(sessions []*session.Session, c Criteria) []*session.Session {
	var result []*session.Session
	for _, s := range sessions {
		if c.UserName != "" && !strings.EqualFold(s.UserName, c.UserName) {
			continue
		}
		if c.RemoteIP != "" && s.RemoteIP != c.RemoteIP {
			continue
		}
		if !c.Since.IsZero() && s.StartedAt.Before(c.Since) {
			continue
		}
		if !c.Until.IsZero() && s.StartedAt.After(c.Until) {
			continue
		}
		if c.MinEvents > 0 && len(s.Events) < c.MinEvents {
			continue
		}
		result = append(result, s)
	}
	return result
}
