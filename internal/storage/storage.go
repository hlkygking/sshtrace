package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"sshtrace/internal/session"
)

// Store handles persisting SSH sessions to disk.
type Store struct {
	baseDir string
}

// New creates a new Store that writes session files to baseDir.
func New(baseDir string) (*Store, error) {
	if err := os.MkdirAll(baseDir, 0o750); err != nil {
		return nil, fmt.Errorf("storage: create base dir: %w", err)
	}
	return &Store{baseDir: baseDir}, nil
}

// Save serialises a session to a JSON file named by session ID and date.
func (s *Store) Save(sess *session.Session) error {
	data, err := json.MarshalIndent(sess, "", "  ")
	if err != nil {
		return fmt.Errorf("storage: marshal session: %w", err)
	}

	filename := fmt.Sprintf("%s_%s.json",
		sess.StartedAt.UTC().Format("20060102T150405Z"),
		sess.ID,
	)
	path := filepath.Join(s.baseDir, filename)

	if err := os.WriteFile(path, data, 0o640); err != nil {
		return fmt.Errorf("storage: write file %s: %w", path, err)
	}
	return nil
}

// Load reads a session from disk by its file path.
func (s *Store) Load(path string) (*session.Session, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("storage: read file %s: %w", path, err)
	}

	var sess session.Session
	if err := json.Unmarshal(data, &sess); err != nil {
		return nil, fmt.Errorf("storage: unmarshal session: %w", err)
	}
	return &sess, nil
}

// ListSince returns file paths for all sessions saved on or after the given time.
func (s *Store) ListSince(since time.Time) ([]string, error) {
	entries, err := os.ReadDir(s.baseDir)
	if err != nil {
		return nil, fmt.Errorf("storage: read dir: %w", err)
	}

	var paths []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		if !info.ModTime().Before(since) {
			paths = append(paths, filepath.Join(s.baseDir, e.Name()))
		}
	}
	return paths, nil
}
