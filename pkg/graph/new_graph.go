//ff:func feature=graph type=helper control=sequence
//ff:what NewGraph — 빈 reins Graph를 생성한다. 내부 *toulmin.Graph 1개와 빈 노드 목록·supersession 맵을 초기화한다. Warrant/Counter로 노드를 더하고 Evaluate로 판독한다.

package graph

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// NewGraph creates an empty reins Graph wrapping a fresh toulmin defeats graph.
// Add nodes with Warrant/Counter and read the verdict with Evaluate.
func NewGraph(name string) *Graph {
	return &Graph{
		tg:         toulmin.NewGraph(name),
		supersedes: make(map[string]map[string]bool),
	}
}
