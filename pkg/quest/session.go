package quest

import (
	"encoding/json"
	"fmt"
	"os"
)

// Session is the persisted runtime state of one quest run: the work-list plus a
// schema version. It is the single source of truth a disposable agent resumes from.
type Session struct {
	Version int     `json:"version"`
	Items   []*Item `json:"items"`
}

// New returns an empty session at the current schema version.
func New() *Session { return &Session{Version: 1} }

// Load reads a session from a JSON file. A missing file is reported via
// os.IsNotExist(err) so callers can create a new session.
func Load(path string) (*Session, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var s Session
	if err := json.Unmarshal(b, &s); err != nil {
		return nil, fmt.Errorf("session %s: %w", path, err)
	}
	return &s, nil
}

// Save writes the session to a JSON file (pretty-printed for diff-friendliness).
func (s *Session) Save(path string) error {
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}

// NextTODO returns the first item still in TODO, or nil when none remain.
func (s *Session) NextTODO() *Item {
	for _, it := range s.Items {
		if it.State == TODO {
			return it
		}
	}
	return nil
}

// Find returns the item with the given Key, or an error if absent.
func (s *Session) Find(key string) (*Item, error) {
	for _, it := range s.Items {
		if it.Key == key {
			return it, nil
		}
	}
	return nil, fmt.Errorf("item not found: %s", key)
}

// Progress tallies items by state (for the status command).
func (s *Session) Progress() map[State]int {
	out := map[State]int{}
	for _, it := range s.Items {
		out[it.State]++
	}
	return out
}
