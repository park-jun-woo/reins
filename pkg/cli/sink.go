package cli

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/park-jun-woo/reins/pkg/quest"
)

// jsonlSink is a quest.Sink that appends one JSON line per item to a file (JSONL).
// It is the default export sink; export is incremental and idempotent because
// quest.Export only emits terminal, not-yet-emitted items.
type jsonlSink struct {
	path string
}

// newJSONLSink returns a sink writing to path, creating the parent directory.
func newJSONLSink(path string) (*jsonlSink, error) {
	if dir := filepath.Dir(path); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, err
		}
	}
	return &jsonlSink{path: path}, nil
}

// Emit appends one JSON-encoded item, newline-terminated, to the sink's file.
func (s *jsonlSink) Emit(it *quest.Item) error {
	f, err := os.OpenFile(s.path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	b, err := json.Marshal(it)
	if err != nil {
		return err
	}
	b = append(b, '\n')
	_, err = f.Write(b)
	return err
}

var _ quest.Sink = (*jsonlSink)(nil)
