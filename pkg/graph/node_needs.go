//ff:func feature=graph type=helper control=sequence
//ff:what Needs — 이 카운터가 읽는 네트워크 ground 의존을 명시 선언한다(열린결정 #1). 인자는 ground 이름들(예 "source-body")로, staged 평가의 tier 분류 입력이 된다: Needs 비어있음=tier0(무네트워크), 하나라도 선언=tier1(ground 의존). tier0가 잔존 Fail이면 tier1을 평가하지 않으므로 그 ground는 resolve되지 않는다(G5). 체이닝을 위해 이 노드를 반환한다.

package graph

// Needs declares the network ground dependencies this counter reads (open decision
// #1). The arguments are ground names (e.g. "source-body"); they are the input to
// the staged evaluator's tier classification: no Needs ⇒ tier-0 (no network), any
// Needs ⇒ tier-1 (ground-dependent). When tier-0 has a residual Fail, tier-1 is not
// evaluated and so its grounds are never resolved (G5). Returns this node for
// chaining.
func (n *Node) Needs(grounds ...string) *Node {
	n.needs = append(n.needs, grounds...)
	return n
}
