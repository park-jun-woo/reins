//ff:func feature=graph type=helper control=sequence
//ff:what Attacks — toulmin 그래프에 defeat 엣지를 배선한다(이 노드 → target, 보통 counter→warrant). verdict 계산·trace 기록·Activated 판독의 토대. staged 평가가 tier 서브그래프를 재배선할 수 있게 target을 reins-side에도 기록한다. 체이닝을 위해 이 노드를 반환한다.

package graph

// Attacks wires a toulmin defeat edge from this node to target (typically
// counter→warrant). This is the edge toulmin uses for verdict computation and
// trace recording. Returns this node for chaining.
func (n *Node) Attacks(target *Node) *Node {
	n.tr.Attacks(target.tr)
	n.attacks = append(n.attacks, target)
	return n
}
