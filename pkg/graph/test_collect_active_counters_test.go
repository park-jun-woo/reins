//ff:func feature=graph type=helper control=sequence
//ff:what collectActiveCounters 단위증명 — 워런트 trace를 훑어 Activated 카운터만 (노드,Fact)로, 등록 순서대로, 중복 제거해 수집하는지; 비활성·워런트·미매칭 엔트리를 제외하고 여러 워런트에 걸친 중복 노드를 한 번만 담는지 검증한다.

package graph

import (
	"testing"

	"github.com/park-jun-woo/reins/pkg/quest"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestCollectActiveCounters(t *testing.T) {
	g := NewGraph("collect")
	// register nodes manually (no toulmin wiring needed for this unit)
	warN := &Node{graph: g, id: "w", shortName: "w", isWarrant: true}
	a := &Node{graph: g, id: "a", shortName: "a"}
	b := &Node{graph: g, id: "b", shortName: "b"}
	c := &Node{graph: g, id: "c", shortName: "c"}
	g.nodes = []*Node{warN, a, b, c}

	te := func(name string, act bool, ev any) toulmin.TraceEntry {
		return toulmin.TraceEntry{Name: name, Activated: act, Evidence: ev}
	}

	results := []toulmin.EvalResult{
		{Name: "w", Trace: []toulmin.TraceEntry{
			te("w", true, nil),                          // warrant -> skipped
			te("a", true, quest.Fact{Where: "wa"}),      // active counter a
			te("b", false, quest.Fact{Where: "wb"}),     // inactive -> skipped
			te("zzz", true, quest.Fact{Where: "ghost"}), // unmatched -> skipped
		}},
		{Name: "w2", Trace: []toulmin.TraceEntry{
			te("a", true, quest.Fact{Where: "dup"}), // duplicate a -> first wins
			te("c", true, quest.Fact{Where: "wc"}),  // active counter c
		}},
	}

	got := g.collectActiveCounters(results)

	// Expect a and c only, in registration order (a before c).
	if len(got) != 2 {
		t.Fatalf("got %d active counters, want 2: %+v", len(got), got)
	}
	if got[0].node.id != "a" || got[1].node.id != "c" {
		t.Fatalf("order wrong: got %s,%s want a,c", got[0].node.id, got[1].node.id)
	}
	if got[0].fact.Where != "wa" {
		t.Fatalf("a fact = %q want wa (first occurrence wins)", got[0].fact.Where)
	}
	if got[1].fact.Where != "wc" {
		t.Fatalf("c fact = %q want wc", got[1].fact.Where)
	}
}
