//ff:func feature=graph type=helper control=iteration dimension=2
//ff:what collectActiveCounters — 워런트별 trace를 훑어 Activated인 카운터 노드를 (노드, Fact)로 수집한다. 엔트리 1건 기록은 recordActiveCounter에 위임(중복 제거·Fact 추출). 등록 순서를 보존해 결정론을 유지하고 워런트 자신은 제외한다.

package graph

import (
	"github.com/park-jun-woo/reins/pkg/quest"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

// collectActiveCounters scans each warrant's trace and collects every activated
// counter node as (node, Fact). Recording a single entry (dedup + Fact extraction)
// is delegated to recordActiveCounter. It preserves node registration order for
// determinism and skips warrant nodes.
func (g *Graph) collectActiveCounters(results []toulmin.EvalResult) []activeCounter {
	facts := make(map[string]quest.Fact)
	seen := make(map[string]bool)
	for _, res := range results {
		for _, te := range res.Trace {
			g.recordActiveCounter(te, seen, facts)
		}
	}
	out := make([]activeCounter, 0, len(seen))
	for _, n := range g.nodes {
		if n.isWarrant || !seen[n.id] {
			continue
		}
		out = append(out, activeCounter{node: n, fact: facts[n.id]})
	}
	return out
}
