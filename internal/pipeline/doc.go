// Package pipeline provides a composable processing chain for SSH sessions.
//
// A Pipeline accepts any number of Processor implementations and applies
// them sequentially to a session. This allows features such as sanitisation,
// redaction, tagging, and truncation to be composed without coupling them
// to each other.
//
// Example usage:
//
//	san, _ := sanitize.New(sanitize.DefaultPatterns())
//	red, _ := redact.New(redact.DefaultFields())
//	p, _ := pipeline.New(san, red)
//	result, err := p.Run(sess)
package pipeline
