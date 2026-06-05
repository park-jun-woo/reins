//ff:func feature=graph type=helper control=sequence
//ff:what PASS 공략집 공백 증명 — 카운터가 발동하지 않아 잔존이 0이면 Evaluate가 PASS를 내고 Verdict.Feedback이 빈 문자열인지(공략집은 FAIL/REVIEW 전용) 검증한다.

package graph

import (
	"testing"

	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
)

func TestRenderWalkthroughPassEmpty(t *testing.T) {
	rA := fireRule("rule-a", gate.LevelFail)
	g := NewGraph("pass")
	pass := g.Warrant(alwaysTrue)
	g.Counter(rA, gate.LevelFail).Attacks(pass)

	v := g.Evaluate(gate.Context{Submission: map[string]bool{}})
	if v.Outcome != quest.OutPass {
		t.Fatalf("outcome=%s want PASS", v.Outcome)
	}
	if v.Feedback != "" {
		t.Fatalf("PASS Feedback = %q, want empty", v.Feedback)
	}
}
