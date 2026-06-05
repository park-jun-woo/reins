//ff:func feature=graph type=helper control=sequence
//ff:what collectActiveCounters non-Fact 증명 — Activated이지만 Evidence가 quest.Fact가 아니면 노드는 수집하되 Fact는 제로값으로 남는지 검증한다.

package graph

import (
	"testing"

	"github.com/park-jun-woo/reins/pkg/quest"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestCollectActiveCountersNonFactEvidence(t *testing.T) {
	g := NewGraph("collect2")
	a := &Node{graph: g, id: "a", shortName: "a"}
	g.nodes = []*Node{a}
	// Activated but evidence is not a quest.Fact -> node still collected, zero Fact.
	results := []toulmin.EvalResult{
		{Name: "w", Trace: []toulmin.TraceEntry{
			{Name: "a", Activated: true, Evidence: 42},
		}},
	}
	got := g.collectActiveCounters(results)
	if len(got) != 1 || got[0].node.id != "a" {
		t.Fatalf("got %+v want single node a", got)
	}
	if (got[0].fact != quest.Fact{}) {
		t.Fatalf("fact = %+v want zero value", got[0].fact)
	}
}
