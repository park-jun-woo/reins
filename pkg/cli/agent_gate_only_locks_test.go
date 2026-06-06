//ff:func feature=cli type=command control=sequence level=error
//ff:what TestAgentGateOnlyLocks — backend가 무엇을 뱉든 게이트가 실패시키면 PASS에 닿지 않고 MaxTries 소진 후 DONE으로 종료(잠금 권한은 게이트만)하는지 검증.

package cli

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/reins/pkg/llm"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// TestAgentGateOnlyLocks: no matter what the backend emits, a submission the gate
// fails never reaches PASS; MaxTries is exhausted and the item is DONE (loop ends).
func TestAgentGateOnlyLocks(t *testing.T) {
	dir := t.TempDir()
	session := dir + "/session.json"
	out := dir + "/out.jsonl"

	calls := 0
	backend := llm.CallFunc(func(system, user string) (string, error) {
		calls++
		return "bad", nil // always fails the gate
	})
	opts := Options{Agent: &AgentOptions{LLM: backend}}

	if _, err := newAgentRoot(t, opts, session, out, "scan", "a"); err != nil {
		t.Fatalf("scan: %v", err)
	}
	got, err := newAgentRoot(t, opts, session, out, "agent")
	if err != nil {
		t.Fatalf("agent: %v", err)
	}
	if calls != quest.MaxTries {
		t.Fatalf("backend called %d times, want MaxTries=%d", calls, quest.MaxTries)
	}
	if strings.Contains(got, "a -> PASS") {
		t.Fatalf("agent must not PASS a gate-failing submission: %q", got)
	}
	// The item must be terminal (DONE) so NextTODO drops it (monotone convergence).
	if !strings.Contains(got, "processed 1 item") {
		t.Fatalf("agent output = %q, want 'processed 1 item'", got)
	}
}
