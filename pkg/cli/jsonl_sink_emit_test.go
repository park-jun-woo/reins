//ff:func feature=cli type=helper control=sequence level=error
//ff:what Emit이 아이템 하나당 JSON 한 줄(개행 종단)을 append하는지 검증한다.

package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/park-jun-woo/reins/pkg/quest"
)

// TestEmitAppendsJSONL: Emit appends one JSON line per item (newline-terminated).
func TestEmitAppendsJSONL(t *testing.T) {
	path := filepath.Join(t.TempDir(), "out.jsonl")
	sink, err := newJSONLSink(path)
	if err != nil {
		t.Fatalf("newJSONLSink: %v", err)
	}
	if err := sink.Emit(&quest.Item{Key: "a", State: quest.PASS}); err != nil {
		t.Fatalf("Emit a: %v", err)
	}
	if err := sink.Emit(&quest.Item{Key: "b", State: quest.DONE}); err != nil {
		t.Fatalf("Emit b: %v", err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	lines := strings.Split(strings.TrimRight(string(b), "\n"), "\n")
	if len(lines) != 2 {
		t.Fatalf("got %d lines, want 2: %q", len(lines), b)
	}
	if !strings.Contains(lines[0], `"key":"a"`) || !strings.Contains(lines[1], `"key":"b"`) {
		t.Fatalf("lines = %q", lines)
	}
}
