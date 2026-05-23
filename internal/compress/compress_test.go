package compress

import (
	"compress/gzip"
	"bytes"
	"testing"
)

func TestNewDefaultLevel(t *testing.T) {
	c, err := New(gzip.DefaultCompression)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil Compressor")
	}
}

func TestNewInvalidLevel(t *testing.T) {
	_, err := New(99)
	if err == nil {
		t.Fatal("expected error for invalid level, got nil")
	}
}

func TestCompressDecompressRoundtrip(t *testing.T) {
	c, _ := New(gzip.DefaultCompression)
	orig := []byte("ssh session data: ls -la && cat /etc/passwd")

	compressed, err := c.Compress(orig)
	if err != nil {
		t.Fatalf("Compress: %v", err)
	}
	if bytes.Equal(compressed, orig) {
		t.Fatal("compressed output should differ from original")
	}

	decompressed, err := c.Decompress(compressed)
	if err != nil {
		t.Fatalf("Decompress: %v", err)
	}
	if !bytes.Equal(decompressed, orig) {
		t.Fatalf("roundtrip mismatch: got %q, want %q", decompressed, orig)
	}
}

func TestCompressReducesSize(t *testing.T) {
	c, _ := New(gzip.BestCompression)
	// Repetitive data compresses well.
	orig := bytes.Repeat([]byte("audit log entry\n"), 100)

	compressed, err := c.Compress(orig)
	if err != nil {
		t.Fatalf("Compress: %v", err)
	}
	if len(compressed) >= len(orig) {
		t.Errorf("expected compressed size < original; got %d >= %d", len(compressed), len(orig))
	}
}

func TestDecompressInvalidData(t *testing.T) {
	c, _ := New(gzip.DefaultCompression)
	_, err := c.Decompress([]byte("not gzip data"))
	if err == nil {
		t.Fatal("expected error decompressing invalid data, got nil")
	}
}

func TestCompressEmptyInput(t *testing.T) {
	c, _ := New(gzip.BestSpeed)
	compressed, err := c.Compress([]byte{})
	if err != nil {
		t.Fatalf("Compress empty: %v", err)
	}
	decompressed, err := c.Decompress(compressed)
	if err != nil {
		t.Fatalf("Decompress empty: %v", err)
	}
	if len(decompressed) != 0 {
		t.Errorf("expected empty output, got %d bytes", len(decompressed))
	}
}
