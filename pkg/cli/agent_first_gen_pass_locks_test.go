//ff:func feature=cli type=command control=sequence level=error
//ff:what TestAgentFirstGenPassLocks — 통과 페이로드를 반환하는 backend가 아이템을 한 번에 PASS로 잠그는지(backend 1회 호출) 무네트워크 검증.

package cli

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/reins/pkg/llm"
)

// TestAgentFirstGenPassLocks: a backend that returns a passing payload locks the
// item PASS in one shot.
func TestAgentFirstGenPassLocks(t *testing.T) {
	dir := t.TempDir()
	session := dir + "/session.json"
	out := dir + "/out.jsonl"

	calls := 0
	backend := llm.CallFunc(func(system, user string) (string, error) {
		calls++
		return "good", nil
	})
	opts := Options{Agent: &AgentOptions{LLM: backend}}

	if _, err := newAgentRoot(t, opts, session, out, "scan", "a"); err != nil {
		t.Fatalf("scan: %v", err)
	}
	got, err := newAgentRoot(t, opts, session, out, "agent")
	if err != nil {
		t.Fatalf("agent: %v", err)
	}
	if !strings.Contains(got, "a -> PASS") {
		t.Fatalf("agent output = %q, want PASS", got)
	}
	if calls != 1 {
		t.Fatalf("backend called %d times, want 1 (PASS in one shot)", calls)
	}
}
