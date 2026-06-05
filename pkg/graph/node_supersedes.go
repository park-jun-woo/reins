//ff:func feature=graph type=helper control=iteration dimension=1
//ff:what Supersedes — reins-side 우선순위(흡수) 관계를 등록한다(toulmin 엣지 아님). 이 노드(상류)가 활성이면 targets(하류) 카운터를 Evaluate 집계에서 제외한다. precedence는 toulmin Attacks로 표현하면 Activated를 못 꺼 곁가지가 FAIL을 남기므로, reins-side 맵으로 두고 읽기 레이어가 크리스프하게 적용한다. 체이닝을 위해 이 노드를 반환한다.

package graph

// Supersedes registers a reins-side precedence (absorption) relation — NOT a
// toulmin edge. When this node (upstream) is active, the target counters
// (downstream) are excluded from Evaluate's aggregation. Expressing precedence as
// a toulmin Attacks edge would not clear the target's Activated flag (so a
// side-branch would still leave a FAIL); hence it lives in a reins-side map that
// the read layer applies crisply. Returns this node for chaining.
func (n *Node) Supersedes(targets ...*Node) *Node {
	set := n.graph.supersedes[n.id]
	if set == nil {
		set = make(map[string]bool)
		n.graph.supersedes[n.id] = set
	}
	for _, t := range targets {
		set[t.id] = true
	}
	return n
}
