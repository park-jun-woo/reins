//ff:func feature=graph type=helper control=sequence
//ff:what Warrant — 워런트 노드를 그래프에 추가한다(tautology PASS 워런트 등). fn은 toulmin 규칙 함수(예: alwaysTrue). toulmin g.Rule(fn).With(고유 idSpec)로 등록해 ruleID를 분리하고, Node(isWarrant=true)를 보관·반환한다. 워런트는 Level 메타가 없다(판정 근거는 카운터 Level).

package graph

import (
	"strconv"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

// Warrant adds a warrant node (e.g. the tautology PASS warrant) to the graph.
// fn is a toulmin rule function such as alwaysTrue. It is registered via
// g.tg.Rule(fn).With(unique idSpec) to keep the ruleID distinct, and the Node is
// recorded and returned. A warrant carries no Level meta — the verdict basis is
// the counters' Levels.
func (g *Graph) Warrant(fn func(toulmin.Context, toulmin.Specs) (bool, any)) *Node {
	spec := idSpec{ID: "warrant#" + strconv.Itoa(len(g.nodes))}
	id := ruleIDFor(fn, spec)
	tr := g.tg.Rule(fn).With(spec)
	n := &Node{
		graph:     g,
		tr:        tr,
		id:        id,
		shortName: shortNameFor(id),
		isWarrant: true,
		fn:        fn,
		spec:      spec,
	}
	g.nodes = append(g.nodes, n)
	return n
}
