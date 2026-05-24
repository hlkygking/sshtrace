package audit

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// ReadOptions controls which records are returned by Read.
type ReadOptions struct {
	Kind      EventKind // filter by event kind; empty means all
	User      string    // filter by user; empty means all
	Since     time.Time // only records at or after this time
}

// Read parses newline-delimited JSON audit records from r and returns those
// that match opts. A zero-value ReadOptions returns all records.
func Read(r io.Reader, opts ReadOptions) ([]Record, error) {
	var records []Record
	scanner := bufio.NewScanner(r)
	line := 0
	for scanner.Scan() {
		line++
		var rec Record
		if err := json.Unmarshal(scanner.Bytes(), &rec); err != nil {
			return nil, fmt.Errorf("audit: parse line %d: %w", line, err)
		}
		if opts.Kind != "" && rec.Kind != opts.Kind {
			continue
		}
		if opts.User != "" && rec.User != opts.User {
			continue
		}
		if !opts.Since.IsZero() && rec.Timestamp.Before(opts.Since) {
			continue
		}
		records = append(records, rec)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("audit: scan: %w", err)
	}
	return records, nil
}
