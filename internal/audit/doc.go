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
// Example:
//
//	logger := audit.New(os.Stdout)
//	logger.LogSession(sess, audit.EventSessionStarted)
package audit
