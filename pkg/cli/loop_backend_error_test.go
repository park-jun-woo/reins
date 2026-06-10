//ff:func feature=cli type=command control=sequence level=error
//ff:what TestLoopBackendError — backend 에러가 루프를 그 에러로 중단시키는지 검증.

package cli

import (
	"testing"

	"github.com/park-jun-woo/reins/pkg/llm"
)

// TestLoopBackendError: a backend error aborts the loop with that error.
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
	if _, err := newLoopRoot(t, opts, session, out, "loop"); err == nil {
		t.Fatal("loop = nil error, want backend error")
	}
}
