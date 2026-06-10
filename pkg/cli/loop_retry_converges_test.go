//ff:func feature=cli type=command control=sequence level=error
//ff:what TestLoopRetryConverges — 첫 시도 FAIL 후 재시도 PASS로 루프가 수렴하는지(backend 2회 호출, FAIL·PASS 출력) 검증.

package cli

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/reins/pkg/llm"
)

// TestLoopRetryConverges: FAIL then PASS — the loop retries and converges.
func TestLoopRetryConverges(t *testing.T) {
	dir := t.TempDir()
	session := dir + "/session.json"
	out := dir + "/out.jsonl"

	calls := 0
	backend := llm.CallFunc(func(system, user string) (string, error) {
		calls++
		if calls == 1 {
			return "bad", nil // first attempt fails the gate
		}
		return "good", nil // retry passes
	})
	opts := Options{Loop: &LoopOptions{LLM: backend}}

	if _, err := newLoopRoot(t, opts, session, out, "scan", "a"); err != nil {
		t.Fatalf("scan: %v", err)
	}
	got, err := newLoopRoot(t, opts, session, out, "loop")
	if err != nil {
		t.Fatalf("loop: %v", err)
	}
	if calls != 2 {
		t.Fatalf("backend called %d times, want 2 (FAIL then PASS)", calls)
	}
	if !strings.Contains(got, "a -> FAIL") || !strings.Contains(got, "a -> PASS") {
		t.Fatalf("loop output = %q, want a FAIL then a PASS", got)
	}
}
