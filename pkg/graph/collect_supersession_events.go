//ff:func feature=graph type=helper control=iteration dimension=2
//ff:what collectSupersessionEvents — 활성 카운터 집합에서 흡수된 곁가지 이벤트를 (등록순) 수집한다. 각 활성 상류가 흡수하는 하류 ID 중 그 자신도 활성인 것만 supersessionEvent(absorbed=하류, by=상류)로 기록한다. 같은 하류를 흡수하는 상류가 여럿이면 등록순 첫 상류를 by로 택한다(결정론). applySupersession의 제외 판정과 동일 모집단 — 렌더 전용 부산물.

package graph

// collectSupersessionEvents collects the absorbed side-branch events (in
// registration order) from the active counter set. For each active upstream, every
// downstream it absorbs that is itself active is recorded as a supersessionEvent
// (absorbed=downstream, by=upstream). If several active upstreams absorb the same
// downstream, the first upstream in registration order is chosen as by
// (deterministic). It is the render-only by-product of the same exclusion the read
// layer applies in applySupersession.
func (g *Graph) collectSupersessionEvents(active []activeCounter) []supersessionEvent {
	activeAC := make(map[string]activeCounter, len(active))
	for _, ac := range active {
		activeAC[ac.node.id] = ac
	}
	byOf := make(map[string]activeCounter) // absorbed id -> chosen superseder (first in order)
	seen := make(map[string]bool)
	for _, ac := range active { // active list is registration order
		for downID := range g.supersedes[ac.node.id] {
			if _, ok := activeAC[downID]; !ok || seen[downID] {
				continue
			}
			seen[downID] = true
			byOf[downID] = ac
		}
	}
	out := make([]supersessionEvent, 0, len(byOf))
	for _, n := range g.nodes { // emit in registration order for determinism
		if by, ok := byOf[n.id]; ok {
			out = append(out, supersessionEvent{absorbed: activeAC[n.id], by: by})
		}
	}
	return out
}
