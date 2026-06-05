//ff:func feature=graph type=helper control=sequence
//ff:what Evaluate 엔진오류 분기증명 — 그래프에 사이클이 있으면 결정적 FAIL Fact(Rule="graph", Where="engine")로 환원하는지 검증한다.

package graph

import (
	"testing"

	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
)

func TestEvaluateEngineErrorIsFail(t *testing.T) {
	rFmt := fireRule("fmt", gate.LevelFail)

	g := NewGraph("cycle")
	pass := g.Warrant(alwaysTrue)
	c := g.Counter(rFmt, gate.LevelFail).Attacks(pass)
	// Form a 2-cycle: warrant attacks the counter back.
	pass.Attacks(c)

	v := g.Evaluate(gate.Context{Submission: map[string]bool{"fmt": true}})
	if v.Outcome != quest.OutFail {
		t.Fatalf("engine error: outcome=%s want FAIL", v.Outcome)
	}
	if len(v.Facts) != 1 {
		t.Fatalf("engine error: expected 1 fact, got %d", len(v.Facts))
	}
	if v.Facts[0].Rule != "graph" || v.Facts[0].Where != "engine" {
		t.Fatalf("engine error fact = %+v want {Rule:graph Where:engine}", v.Facts[0])
	}
	if v.Facts[0].Actual == "" {
		t.Fatalf("engine error fact.Actual is empty; want the engine error message")
	}
}
