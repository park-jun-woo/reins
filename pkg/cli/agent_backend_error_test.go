//ff:func feature=cli type=command control=sequence level=error
//ff:what TestAgentBackendError — backend 에러가 루프를 그 에러로 중단시키는지 검증.

package cli

import (
	"testing"

	"github.com/park-jun-woo/reins/pkg/llm"
)

// TestAgentBackendError: a backend error aborts the loop with that error.
func TestAgentBackendError(t *testing.T) {
	dir := t.TempDir()
	session := dir + "/session.json"
	out := dir + "/out.jsonl"

	backend := llm.CallFunc(func(system, user string) (string, error) {
		return "", errBackend
	})
	opts := Options{Agent: &AgentOptions{LLM: backend}}

	if _, err := newAgentRoot(t, opts, session, out, "scan", "a"); err != nil {
		t.Fatalf("scan: %v", err)
	}
	if _, err := newAgentRoot(t, opts, session, out, "agent"); err == nil {
		t.Fatal("agent = nil error, want backend error")
	}
}
