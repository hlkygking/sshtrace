// Package compress provides gzip compression and decompression
// for serialized SSH session data before storage or export.
package compress

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
)

// Compressor handles gzip compression and decompression of raw bytes.
type Compressor struct {
	level int
}

// New returns a Compressor using the given gzip compression level.
// Use compress/gzip constants (e.g. gzip.BestSpeed, gzip.BestCompression,
// gzip.DefaultCompression).
func New(level int) (*Compressor, error) {
	if level != gzip.DefaultCompression && (level < gzip.BestSpeed || level > gzip.BestCompression) {
		return nil, fmt.Errorf("compress: invalid level %d", level)
	}
	return &Compressor{level: level}, nil
}

// Compress compresses src using gzip and returns the compressed bytes.
func (c *Compressor) Compress(src []byte) ([]byte, error) {
	var buf bytes.Buffer
	w, err := gzip.NewWriterLevel(&buf, c.level)
	if err != nil {
		return nil, fmt.Errorf("compress: create writer: %w", err)
	}
	if _, err := w.Write(src); err != nil {
		return nil, fmt.Errorf("compress: write: %w", err)
	}
	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("compress: close writer: %w", err)
	}
	return buf.Bytes(), nil
}

// Decompress decompresses gzip-encoded src and returns the original bytes.
func (c *Compressor) Decompress(src []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(src))
	if err != nil {
		return nil, fmt.Errorf("compress: create reader: %w", err)
	}
	defer r.Close()
	out, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("compress: read: %w", err)
	}
	return out, nil
}
