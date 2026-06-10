//ff:func feature=cli type=command control=sequence level=error
//ff:what TestLoopFirstGenPassLocks — 통과 페이로드를 반환하는 backend가 아이템을 한 번에 PASS로 잠그는지(backend 1회 호출) 무네트워크 검증.

package cli

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/reins/pkg/llm"
)

// TestLoopFirstGenPassLocks: a backend that returns a passing payload locks the
// item PASS in one shot.
func TestLoopFirstGenPassLocks(t *testing.T) {
	dir := t.TempDir()
	session := dir + "/session.json"
	out := dir + "/out.jsonl"

	calls := 0
	backend := llm.CallFunc(func(system, user string) (string, error) {
		calls++
		return "good", nil
	})
	opts := Options{Loop: &LoopOptions{LLM: backend}}

	if _, err := newLoopRoot(t, opts, session, out, "scan", "a"); err != nil {
		t.Fatalf("scan: %v", err)
	}
	got, err := newLoopRoot(t, opts, session, out, "loop")
	if err != nil {
		t.Fatalf("loop: %v", err)
	}
	if !strings.Contains(got, "a -> PASS") {
		t.Fatalf("loop output = %q, want PASS", got)
	}
	if calls != 1 {
		t.Fatalf("backend called %d times, want 1 (PASS in one shot)", calls)
	}
}
