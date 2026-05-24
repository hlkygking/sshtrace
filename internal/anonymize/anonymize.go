// Package anonymize replaces personally identifiable information in SSH
// session metadata (username, remote IP) with stable pseudonyms so that
// sessions can be shared or analysed without exposing real identities.
package anonymize

import (
	"crypto/sha256"
	"fmt"
	"sync"

	"github.com/sshtrace/sshtrace/internal/session"
)

// Anonymizer replaces user and IP fields with deterministic pseudonyms.
type Anonymizer struct {
	mu      sync.Mutex
	salt    string
	userMap map[string]string
	ipMap   map[string]string
}

// New returns an Anonymizer seeded with the provided salt.
// The same salt always produces the same pseudonyms, allowing
// cross-session correlation without revealing real values.
func New(salt string) (*Anonymizer, error) {
	if salt == "" {
		return nil, fmt.Errorf("anonymize: salt must not be empty")
	}
	return &Anonymizer{
		salt:    salt,
		userMap: make(map[string]string),
		ipMap:   make(map[string]string),
	}, nil
}

// Apply returns a shallow copy of s with User and RemoteIP replaced by
// stable pseudonyms. The original session is never modified.
func (a *Anonymizer) Apply(s *session.Session) (*session.Session, error) {
	if s == nil {
		return nil, fmt.Errorf("anonymize: session must not be nil")
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	copy := *s
	copy.User = a.pseudonym("user", s.User, a.userMap)
	copy.RemoteIP = a.pseudonym("ip", s.RemoteIP, a.ipMap)
	return &copy, nil
}

// pseudonym returns a cached or newly computed pseudonym for value.
func (a *Anonymizer) pseudonym(prefix, value string, cache map[string]string) string {
	if p, ok := cache[value]; ok {
		return p
	}
	h := sha256.Sum256([]byte(a.salt + ":" + prefix + ":" + value))
	p := fmt.Sprintf("%s-%x", prefix, h[:4])
	cache[value] = p
	return p
}
