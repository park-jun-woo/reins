//ff:type feature=graph type=model
//ff:what Graph — toulmin을 감싸는 reins 편의 빌더. 내부 *toulmin.Graph(그래프 수학·trace)와 reins-side 노드 목록·supersession 맵(상류 노드 ruleID → 흡수되는 하류 노드 ruleID 집합)을 보유한다. NewGraph로 생성, Warrant/Counter로 노드 추가, Evaluate로 판독.

package graph

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// Graph is a reins convenience builder wrapping a toulmin defeats graph.
// It holds the underlying *toulmin.Graph (graph math + trace) plus reins-side
// node metadata and a supersession map (upstream ruleID → set of absorbed
// downstream ruleIDs). Supersession is reins-side precedence, not a toulmin edge.
type Graph struct {
	tg         *toulmin.Graph
	nodes      []*Node
	supersedes map[string]map[string]bool
}
