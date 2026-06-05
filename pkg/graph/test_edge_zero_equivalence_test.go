//ff:func feature=graph type=helper control=iteration dimension=2
//ff:what 엣지0 동치 증명 — FromRules(rules).Evaluate(ctx)가 현 gate.Evaluate(rules, ctx)와 모든 입력에서 동일 Outcome·동일 Fact 집합(Where 기준)을 내는지 검증한다. 여러 발동 조합(전무·Fail만·Review만·혼합·전부)을 돌려 상위집합 통찰(현 모델=엣지0 그래프)을 못 박는다.

package graph

import (
	"testing"

	"github.com/park-jun-woo/reins/pkg/gate"
)

func TestEdgeZeroEquivalence(t *testing.T) {
	rules := []gate.Rule{
		fireRule("fmt", gate.LevelFail),
		fireRule("holder", gate.LevelFail),
		fireRule("free", gate.LevelReview),
		fireRule("group", gate.LevelReview),
	}
	g := FromRules(rules)

	cases := []map[string]bool{
		{},
		{"fmt": true},
		{"holder": true},
		{"free": true},
		{"group": true},
		{"free": true, "group": true},
		{"fmt": true, "free": true},
		{"fmt": true, "holder": true, "free": true, "group": true},
		{"holder": true, "group": true},
	}

	for i, fired := range cases {
		ctx := gate.Context{Submission: fired}
		gotGraph := g.Evaluate(ctx)
		gotGate := gate.Evaluate(rules, ctx)

		if gotGraph.Outcome != gotGate.Outcome {
			t.Fatalf("case %d %v: graph outcome=%s, gate outcome=%s",
				i, fired, gotGraph.Outcome, gotGate.Outcome)
		}
		gw, ge := factWheres(gotGraph), factWheres(gotGate)
		if !equalStringSets(gw, ge) {
			t.Fatalf("case %d %v: facts differ\n graph=%v\n gate=%v", i, fired, gw, ge)
		}
	}
}
