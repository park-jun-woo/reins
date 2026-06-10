//ff:func feature=cli type=command control=sequence level=error
//ff:what TestLoopMaxItems — --max-items가 처리 아이템 수를 제한하는지(3개 시드, --max-items 1 → "processed 1 item") 검증.

package cli

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/reins/pkg/llm"
)

// TestLoopMaxItems: --max-items caps the number of items processed.
func TestLoopMaxItems(t *testing.T) {
	dir := t.TempDir()
	session := dir + "/session.json"
	out := dir + "/out.jsonl"

	backend := llm.CallFunc(func(system, user string) (string, error) { return "good", nil })
	opts := Options{Loop: &LoopOptions{LLM: backend}}

	if _, err := newLoopRoot(t, opts, session, out, "scan", "a", "b", "c"); err != nil {
		t.Fatalf("scan: %v", err)
	}
	got, err := newLoopRoot(t, opts, session, out, "loop", "--max-items", "1")
	if err != nil {
		t.Fatalf("loop: %v", err)
	}
	if !strings.Contains(got, "processed 1 item") {
		t.Fatalf("loop output = %q, want 'processed 1 item' with --max-items 1", got)
	}
}
