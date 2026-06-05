//ff:func feature=graph type=helper control=iteration dimension=1
//ff:what FromRules — 후방호환 동치 빌더. tautology PASS 워런트 1개를 만들고 모든 gate.Rule을 Counter로 등록해 그 워런트를 Attacks하게 배선한다. supersession 0개(엣지0 그래프). 이 그래프의 Evaluate는 현 gate.Evaluate(rules, ctx)와 동일 Verdict를 낸다(엣지0 동치 — 동치 테스트로 증명).

package graph

import "github.com/park-jun-woo/reins/pkg/gate"

// FromRules builds the back-compat equivalent graph: one tautology PASS warrant
// plus every gate.Rule registered as a Counter attacking that warrant, with zero
// supersession (an edge-zero graph). This graph's Evaluate yields the same Verdict
// as the current gate.Evaluate(rules, ctx) — the edge-zero equivalence proved by
// the equivalence test.
func FromRules(rules []gate.Rule) *Graph {
	g := NewGraph("from-rules")
	pass := g.Warrant(alwaysTrue)
	for _, r := range rules {
		g.Counter(r, r.Meta.Level).Attacks(pass)
	}
	return g
}
