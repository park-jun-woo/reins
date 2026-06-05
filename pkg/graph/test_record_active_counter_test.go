//ff:func feature=graph type=helper control=sequence
//ff:what recordActiveCounter 단위증명 — trace 엔트리 1건을 판정: 비활성·미매칭(nil)·워런트·이미 본 노드는 무시하고, 활성 카운터는 seen 표시 + Fact를 한 번만 담는지(첫 발생 우선, non-Fact evidence는 미기록) 테이블로 커버한다.

package graph

import (
	"testing"

	"github.com/park-jun-woo/reins/pkg/quest"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestRecordActiveCounter(t *testing.T) {
	g := NewGraph("rec")
	warN := &Node{graph: g, id: "w", shortName: "w", isWarrant: true}
	cN := &Node{graph: g, id: "c", shortName: "c"}
	g.nodes = []*Node{warN, cN}

	t.Run("inactive -> ignored", func(t *testing.T) {
		seen, facts := map[string]bool{}, map[string]quest.Fact{}
		g.recordActiveCounter(toulmin.TraceEntry{Name: "c", Activated: false}, seen, facts)
		if len(seen) != 0 || len(facts) != 0 {
			t.Fatalf("inactive recorded: seen=%v facts=%v", seen, facts)
		}
	})

	t.Run("unmatched name -> ignored", func(t *testing.T) {
		seen, facts := map[string]bool{}, map[string]quest.Fact{}
		g.recordActiveCounter(toulmin.TraceEntry{Name: "ghost", Activated: true}, seen, facts)
		if len(seen) != 0 {
			t.Fatalf("unmatched recorded: %v", seen)
		}
	})

	t.Run("warrant -> ignored", func(t *testing.T) {
		seen, facts := map[string]bool{}, map[string]quest.Fact{}
		g.recordActiveCounter(toulmin.TraceEntry{Name: "w", Activated: true}, seen, facts)
		if seen["w"] {
			t.Fatalf("warrant recorded")
		}
	})

	t.Run("already seen -> no overwrite", func(t *testing.T) {
		seen := map[string]bool{"c": true}
		facts := map[string]quest.Fact{"c": {Where: "first"}}
		g.recordActiveCounter(toulmin.TraceEntry{Name: "c", Activated: true, Evidence: quest.Fact{Where: "second"}}, seen, facts)
		if facts["c"].Where != "first" {
			t.Fatalf("overwrote seen fact: %q want first", facts["c"].Where)
		}
	})

	t.Run("active counter with Fact -> recorded", func(t *testing.T) {
		seen, facts := map[string]bool{}, map[string]quest.Fact{}
		g.recordActiveCounter(toulmin.TraceEntry{Name: "c", Activated: true, Evidence: quest.Fact{Where: "hit"}}, seen, facts)
		if !seen["c"] {
			t.Fatalf("not marked seen")
		}
		if facts["c"].Where != "hit" {
			t.Fatalf("fact = %q want hit", facts["c"].Where)
		}
	})

	t.Run("active counter with non-Fact evidence -> seen, no fact", func(t *testing.T) {
		seen, facts := map[string]bool{}, map[string]quest.Fact{}
		g.recordActiveCounter(toulmin.TraceEntry{Name: "c", Activated: true, Evidence: "nope"}, seen, facts)
		if !seen["c"] {
			t.Fatalf("not marked seen")
		}
		if _, ok := facts["c"]; ok {
			t.Fatalf("recorded non-Fact evidence: %v", facts["c"])
		}
	})
}
