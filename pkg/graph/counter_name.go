//ff:func feature=graph type=helper control=sequence
//ff:what counterName — 카운터의 사람이 읽을 이름을 고른다. Fact.Rule(게이트 규칙 ID 스탬프)이 있으면 그것을, 없으면 노드 shortName으로 폴백한다. 공략집 렌더의 결정타·곁가지 라벨에 쓰인다.

package graph

// counterName returns a human-readable label for a counter: the Fact.Rule (the
// stamped gate rule ID) when present, falling back to the node's short name. It
// labels the root cause and side-branches in the walkthrough.
func counterName(ac activeCounter) string {
	if ac.fact.Rule != "" {
		return ac.fact.Rule
	}
	return ac.node.shortName
}
