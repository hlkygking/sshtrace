// Package export provides utilities for exporting SSH session audit data
// to human-readable or machine-parseable formats.
//
// Supported formats:
//
//	- FormatJSON: structured JSON, suitable for programmatic consumption
//	- FormatText: plain text summary, suitable for human review
//
// Example usage:
//
//	ex := export.New(os.Stdout, export.FormatText)
//	if err := ex.Export(sess); err != nil {
//	    log.Fatal(err)
//	}
package export
