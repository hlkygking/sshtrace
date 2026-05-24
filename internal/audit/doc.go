// Package audit implements structured audit logging for sshtrace.
//
// Each audit record captures the session ID, user, remote IP, event kind,
// and an optional detail string. Records are written as newline-delimited
// JSON to any io.Writer, making them easy to ship to log aggregators or
// append to a local file.
//
// Supported event kinds:
//
//   - EventSessionStarted  – emitted when a new SSH session begins
//   - EventSessionClosed   – emitted when a session ends
//   - EventCommandRun      – emitted when a command is captured
//   - EventAlertTriggered  – emitted when an alert rule fires
//
// Usage:
//
// Create a logger by passing any io.Writer to audit.New. For production use,
// consider wrapping the writer with a mutex or using a buffered writer to
// avoid contention under high session concurrency.
//
// Example:
//
//	logger := audit.New(os.Stdout)
//	logger.LogSession(sess, audit.EventSessionStarted)
//
//	// Log a captured command with detail:
//	logger.LogCommand(sess, "/bin/bash -c 'ls -la'")
//
//	// Log an alert firing:
//	logger.LogAlert(sess, "keyword match: 'passwd'")
package audit
