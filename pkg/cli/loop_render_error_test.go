//ff:func feature=cli type=command control=sequence level=error
//ff:what TestLoopRenderError — 루프 안 def.Render 에러가 그 에러로 중단시키는지 검증(정상 def로 시드 후 render-에러 def로 loop 실행).

package cli

import (
	"testing"

	"github.com/park-jun-woo/reins/pkg/llm"
)

// TestLoopRenderError: a def.Render error inside the loop aborts with that error.
func TestLoopRenderError(t *testing.T) {
	dir := t.TempDir()
	session := dir + "/session.json"
	out := dir + "/out.jsonl"

	backend := llm.CallFunc(func(string, string) (string, error) { return "good", nil })
	opts := Options{Loop: &LoopOptions{LLM: backend}}

	// Seed with a non-erroring def first so the item exists, then run the loop with
	// a render-erroring def over the same session.
	if _, err := newLoopRootDef(t, stubDef{}, opts, session, out, "scan", "a"); err != nil {
		t.Fatalf("scan: %v", err)
	}
	if _, err := newLoopRootDef(t, errDef{renderErr: true}, opts, session, out, "loop"); err == nil {
		t.Fatal("loop = nil error, want Render error")
	}
}
