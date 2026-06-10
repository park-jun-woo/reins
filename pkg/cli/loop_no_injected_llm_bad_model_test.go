//ff:func feature=cli type=command control=sequence level=error
//ff:what TestLoopNoInjectedLLMBadModel — 주입 backend가 없으면 --model을 llm.FromFlag로 해석하고, 잘못된 model flag는 루프 전에 에러내는지 검증.

package cli

import (
	"testing"
)

// TestLoopNoInjectedLLMBadModel: with no injected backend the loop resolves
// --model via llm.FromFlag; an invalid model flag errors before the loop.
func TestLoopNoInjectedLLMBadModel(t *testing.T) {
	dir := t.TempDir()
	session := dir + "/session.json"
	out := dir + "/out.jsonl"

	opts := Options{Loop: &LoopOptions{}} // LLM nil ⇒ FromFlag path
	if _, err := newLoopRoot(t, opts, session, out, "scan", "a"); err != nil {
		t.Fatalf("scan: %v", err)
	}
	// "nocolon" has no backend:model form ⇒ FromFlag returns an error.
	if _, err := newLoopRoot(t, opts, session, out, "loop", "--model", "nocolon"); err == nil {
		t.Fatal("loop = nil error, want FromFlag error for bad --model")
	}
}
