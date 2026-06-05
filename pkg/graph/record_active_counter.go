//ff:func feature=graph type=helper control=selection
//ff:what recordActiveCounter — trace 엔트리 1건을 판정해 활성 카운터면 seen에 표시하고 그 Fact를 facts에 한 번만 담는다. 비활성·미매칭·워런트·이미 본 노드는 무시. collectActiveCounters의 중첩을 낮추는 분리 함수.

package graph

import (
	"github.com/park-jun-woo/reins/pkg/quest"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

// recordActiveCounter inspects one trace entry: if it is an activated, non-warrant
// counter not seen before, it marks seen[id] and records its Fact once. Inactive,
// unmatched, warrant, or already-seen entries are ignored. Extracted to keep
// collectActiveCounters within the nesting-depth limit.
func (g *Graph) recordActiveCounter(te toulmin.TraceEntry, seen map[string]bool, facts map[string]quest.Fact) {
	if !te.Activated {
		return
	}
	n := g.nodeByShortName(te.Name)
	if n == nil || n.isWarrant || seen[n.id] {
		return
	}
	seen[n.id] = true
	if f, ok := te.Evidence.(quest.Fact); ok {
		facts[n.id] = f
	}
}
