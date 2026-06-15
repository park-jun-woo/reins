//ff:func feature=cli type=command control=sequence level=error
//ff:what TestLoopBackendError — backend 에러가 루프를 abort하지 않고 아이템 FAIL로 강등되어 런이 nil로 완주하는지 검증(BUG-002 회귀의 정면 단언).

package cli

import (
	"testing"

	"github.com/park-jun-woo/reins/pkg/llm"
)

// TestLoopBackendError: a backend error no longer aborts the loop — it is demoted to
// a retryable item FAIL, so the run completes without error (BUG-002 regression).
func TestLoopBackendError(t *testing.T) {
	dir := t.TempDir()
	session := dir + "/session.json"
	out := dir + "/out.jsonl"

	backend := llm.CallFunc(func(system, user string) (string, error) {
		return "", errBackend
	})
	opts := Options{Loop: &LoopOptions{LLM: backend}}

	if _, err := newLoopRoot(t, opts, session, out, "scan", "a"); err != nil {
		t.Fatalf("scan: %v", err)
	}
	if _, err := newLoopRoot(t, opts, session, out, "loop"); err != nil {
		t.Fatalf("loop = %v, want nil (backend error is demoted, not propagated)", err)
	}
}
