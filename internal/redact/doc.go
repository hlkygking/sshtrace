// Package redact provides a Redactor that scrubs sensitive field values
// from SSH session event data.
//
// It matches common patterns such as:
//
//	password=secret123
//	token: abc123
//	apikey=xyz
//
// and replaces the value portion with [REDACTED].
//
// Usage:
//
//	r, err := redact.New(nil) // use default fields
//	if err != nil { ... }
//	clean := r.RedactSession(sess)
package redact
