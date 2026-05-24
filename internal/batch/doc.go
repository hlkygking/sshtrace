// Package batch provides bulk-processing support for SSH session audit trails.
//
// A Processor chains one or more Step functions and applies them to a slice
// of sessions in a single call. Each step receives the session and may mutate
// it or return an error to abort further processing for that session.
//
// Steps are applied in the order they are registered. A failure in any step
// stops the remaining steps for that session but does not affect other
// sessions in the batch.
//
// Example usage:
//
//	p, err := batch.New(sanitizeStep, enrichStep, exportStep)
//	if err != nil { ... }
//	results := p.Run(sessions)
//	for _, r := range batch.Errors(results) {
//		log.Printf("session %s failed: %v", r.SessionID, r.Err)
//	}
package batch
