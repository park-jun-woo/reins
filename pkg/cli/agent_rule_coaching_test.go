//ff:func feature=cli type=command control=sequence level=error
//ff:what TestAgentRuleCoaching — FAIL 시 RootCause로 매핑된 RuleSystem 코칭이 다음 시도의 system 프롬프트에 합성되는지(첫 system엔 없고 재시도 system엔 있음) 검증.

package cli

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/reins/pkg/llm"
)

// TestAgentRuleCoaching: on FAIL the RootCause-mapped RuleSystem coaching is
// composed into the next attempt's system prompt.
func TestAgentRuleCoaching(t *testing.T) {
	dir := t.TempDir()
	session := dir + "/session.json"
	out := dir + "/out.jsonl"

	const coach = "COACH-FOR-NOT-BAD"
	var systems []string
	calls := 0
	backend := llm.CallFunc(func(system, user string) (string, error) {
		systems = append(systems, system)
		calls++
		if calls == 1 {
			return "bad", nil
		}
		return "good", nil
	})
	// stubDef's FAIL rule has ID "not-bad" — gate.Evaluate sets RootCause to it.
	opts := Options{Agent: &AgentOptions{
		LLM:        backend,
		RuleSystem: map[string]string{"not-bad": coach},
	}}

	if _, err := newAgentRoot(t, opts, session, out, "scan", "a"); err != nil {
		t.Fatalf("scan: %v", err)
	}
	if _, err := newAgentRoot(t, opts, session, out, "agent"); err != nil {
		t.Fatalf("agent: %v", err)
	}
	if len(systems) < 2 {
		t.Fatalf("expected >=2 attempts, got %d", len(systems))
	}
	if strings.Contains(systems[0], coach) {
		t.Fatalf("first attempt system should not yet carry coaching: %q", systems[0])
	}
	if !strings.Contains(systems[1], coach) {
		t.Fatalf("retry system missing rule coaching: %q", systems[1])
	}
}
