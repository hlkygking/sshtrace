// Package snapshot provides point-in-time capture of session state,
// allowing sessions to be saved and restored at specific moments.
package snapshot

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/sshtrace/sshtrace/internal/session"
)

// Snapshot holds a frozen copy of a session at a specific point in time.
type Snapshot struct {
	SessionID string          `json:"session_id"`
	CapturedAt time.Time      `json:"captured_at"`
	EventCount int            `json:"event_count"`
	Session    *session.Session `json:"session"`
}

// Manager handles saving and loading snapshots to/from disk.
type Manager struct {
	dir string
}

// New creates a new Manager that stores snapshots in dir.
func New(dir string) (*Manager, error) {
	if dir == "" {
		return nil, errors.New("snapshot: directory must not be empty")
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("snapshot: create dir: %w", err)
	}
	return &Manager{dir: dir}, nil
}

// Take creates a Snapshot from the current state of s.
func (m *Manager) Take(s *session.Session) (*Snapshot, error) {
	if s == nil {
		return nil, errors.New("snapshot: session must not be nil")
	}
	snap := &Snapshot{
		SessionID:  s.ID,
		CapturedAt: time.Now().UTC(),
		EventCount: len(s.Events),
		Session:    s,
	}
	path := m.pathFor(s.ID)
	data, err := json.Marshal(snap)
	if err != nil {
		return nil, fmt.Errorf("snapshot: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return nil, fmt.Errorf("snapshot: write: %w", err)
	}
	return snap, nil
}

// Load reads the snapshot for the given session ID from disk.
func (m *Manager) Load(sessionID string) (*Snapshot, error) {
	if sessionID == "" {
		return nil, errors.New("snapshot: session ID must not be empty")
	}
	data, err := os.ReadFile(m.pathFor(sessionID))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("snapshot: not found for session %q", sessionID)
		}
		return nil, fmt.Errorf("snapshot: read: %w", err)
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal: %w", err)
	}
	return &snap, nil
}

// Delete removes the snapshot for the given session ID.
func (m *Manager) Delete(sessionID string) error {
	if sessionID == "" {
		return errors.New("snapshot: session ID must not be empty")
	}
	err := os.Remove(m.pathFor(sessionID))
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("snapshot: delete: %w", err)
	}
	return nil
}

func (m *Manager) pathFor(sessionID string) string {
	return filepath.Join(m.dir, sessionID+".snap.json")
}
