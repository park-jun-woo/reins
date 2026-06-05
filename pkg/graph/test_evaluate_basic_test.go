//ff:func feature=graph type=helper control=sequence
//ff:what 기본 판독 증명 — 워런트만 있는(카운터 발동 0) 그래프는 PASS·Facts 없음을 내고, 활성 Fail 카운터 1개는 FAIL과 그 카운터 Fact(Rule=gate ID 스탬프, Where 전달)를 내는지 검증한다. 어댑터(gate.Context→toulmin ctx)·trace 판독 경로의 스모크 테스트.

package graph

import (
	"testing"

	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
)

func TestEvaluateBasic(t *testing.T) {
	rFmt := fireRule("email-format", gate.LevelFail)

	g := NewGraph("basic")
	pass := g.Warrant(alwaysTrue)
	g.Counter(rFmt, gate.LevelFail).Attacks(pass)

	// No counter fires -> PASS, no facts.
	if v := g.Evaluate(gate.Context{Submission: map[string]bool{}}); v.Outcome != quest.OutPass {
		t.Fatalf("empty: outcome=%s want PASS", v.Outcome)
	} else if len(v.Facts) != 0 {
		t.Fatalf("empty: expected no facts, got %v", v.Facts)
	}

	// Counter fires -> FAIL with its Fact (stamped Rule + Where).
	v := g.Evaluate(gate.Context{Submission: map[string]bool{"email-format": true}})
	if v.Outcome != quest.OutFail {
		t.Fatalf("fired: outcome=%s want FAIL", v.Outcome)
	}
	if len(v.Facts) != 1 {
		t.Fatalf("fired: expected 1 fact, got %d", len(v.Facts))
	}
	if v.Facts[0].Rule != "email-format" {
		t.Fatalf("fired: fact.Rule=%q want email-format", v.Facts[0].Rule)
	}
	if v.Facts[0].Where != "email-format" {
		t.Fatalf("fired: fact.Where=%q want email-format", v.Facts[0].Where)
	}
}
