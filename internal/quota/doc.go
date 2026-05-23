// Package quota provides per-user enforcement of session and event limits
// for the sshtrace audit system.
//
// An Enforcer is created with a Limits configuration that specifies:
//   - MaxSessions: maximum number of recorded sessions per user
//   - MaxEventsPerSession: maximum events allowed within a single session
//   - MaxSessionDuration: maximum wall-clock length of a session
//
// Usage:
//
//	enforcer := quota.New(quota.DefaultLimits())
//
//	if err := enforcer.CheckSession(user); err != nil {
//		// reject new session
//	}
//	enforcer.RecordSession(user)
//
// A zero value for any limit field disables that particular check.
package quota
