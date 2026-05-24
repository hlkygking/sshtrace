// Package summary provides functionality to generate statistical summaries
// of SSH sessions, including command counts, active users, and duration stats.
package summary

import (
	"time"

	"sshtrace/internal/session"
)

// Report holds aggregated statistics for a collection of sessions.
type Report struct {
	TotalSessions  int           `json:"total_sessions"`
	UniqueUsers    []string      `json:"unique_users"`
	UniqueIPs      []string      `json:"unique_ips"`
	TotalEvents    int           `json:"total_events"`
	AvgDuration    time.Duration `json:"avg_duration_ns"`
	LongestSession time.Duration `json:"longest_session_ns"`
	ShortestSession time.Duration `json:"shortest_session_ns"`
}

// UniqueUserCount returns the number of distinct users in the report.
func (r Report) UniqueUserCount() int {
	return len(r.UniqueUsers)
}

// UniqueIPCount returns the number of distinct client IPs in the report.
func (r Report) UniqueIPCount() int {
	return len(r.UniqueIPs)
}

// Generate builds a Report from the provided list of sessions.
func Generate(sessions []*session.Session) Report {
	if len(sessions) == 0 {
		return Report{}
	}

	userSet := map[string]struct{}{}
	ipSet := map[string]struct{}{}
	totalEvents := 0
	var totalDuration time.Duration
	longest := time.Duration(0)
	shortest := time.Duration(-1)

	for _, s := range sessions {
		userSet[s.User] = struct{}{}
		ipSet[s.ClientIP] = struct{}{}
		totalEvents += len(s.Events)

		d := s.Duration()
		totalDuration += d

		if d > longest {
			longest = d
		}
		if shortest < 0 || d < shortest {
			shortest = d
		}
	}

	users := make([]string, 0, len(userSet))
	for u := range userSet {
		users = append(users, u)
	}

	ips := make([]string, 0, len(ipSet))
	for ip := range ipSet {
		ips = append(ips, ip)
	}

	if shortest < 0 {
		shortest = 0
	}

	return Report{
		TotalSessions:   len(sessions),
		UniqueUsers:     users,
		UniqueIPs:       ips,
		TotalEvents:     totalEvents,
		AvgDuration:     totalDuration / time.Duration(len(sessions)),
		LongestSession:  longest,
		ShortestSession: shortest,
	}
}
