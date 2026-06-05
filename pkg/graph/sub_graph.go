//ff:func feature=graph type=helper control=iteration dimension=2
//ff:what subGraph — 원 그래프에서 include(counter)가 참인 카운터 + 모든 워런트만 담은 새 Graph를 만든다. 노드는 cloneNode로 같은 ID로 재등록하고, 원본 Attacks 엣지 중 양끝이 포함된 것만 재배선하며, supersession 엣지도 양끝이 포함된 것만 복사한다. 워런트는 항상 포함(tautology PASS 워런트가 trace의 활성 판독 토대). 결과 그래프의 Evaluate는 기존 판독 경로를 그대로 탄다. staged 평가의 tier0 서브그래프 구성에 쓴다.

package graph

// subGraph builds a new Graph containing every warrant plus the counters for which
// include returns true. Nodes are re-registered with the same IDs via cloneNode,
// the original Attacks edges are rewired only when both endpoints are included, and
// supersession edges are copied only when both endpoints are included. Warrants are
// always included (the tautology PASS warrant is the basis of the trace's
// Activated read). The resulting graph's Evaluate uses the existing read path
// unchanged. It is used to build the tier-0 sub-graph for staged evaluation.
func (g *Graph) subGraph(name string, include func(*Node) bool) *Graph {
	sub := NewGraph(name)
	clone := make(map[string]*Node, len(g.nodes)) // old id → new node
	for _, n := range g.nodes {
		if n.isWarrant || include(n) {
			clone[n.id] = sub.cloneNode(n)
		}
	}
	// Rewire Attacks edges whose both endpoints survived.
	for _, n := range g.nodes {
		src, ok := clone[n.id]
		if !ok {
			continue
		}
		for _, tgt := range n.attacks {
			if dst, ok := clone[tgt.id]; ok {
				src.Attacks(dst)
			}
		}
	}
	// Copy supersession edges among included nodes.
	for upID, downs := range g.supersedes {
		if _, ok := clone[upID]; !ok {
			continue
		}
		for downID := range downs {
			if _, ok := clone[downID]; !ok {
				continue
			}
			set := sub.supersedes[upID]
			if set == nil {
				set = make(map[string]bool)
				sub.supersedes[upID] = set
			}
			set[downID] = true
		}
	}
	return sub
}
