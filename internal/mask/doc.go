// Package mask provides selective field masking for SSH session events.
//
// Use mask.New to create a Masker configured with a list of event-type names
// whose Data field should be replaced with a placeholder string before the
// session is stored, exported, or transmitted.
//
// Example:
//
//	m, err := mask.New([]string{"password", "token"}, "")
//	if err != nil {
//		log.Fatal(err)
//	}
//	masked := m.Apply(sess)
//
// The placeholder defaults to "[MASKED]" when an empty string is provided.
package mask
