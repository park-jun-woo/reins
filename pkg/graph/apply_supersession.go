//ff:func feature=graph type=helper control=iteration dimension=2
//ff:what applySupersession — 활성 카운터 집합에 reins-side Supersedes를 크리스프하게 적용한다. 활성 상류 카운터가 흡수하는 하류 카운터 ID를 제외 집합에 모으고, 활성이면서 제외되지 않은 잔존 카운터만 (등록 순서대로) 반환한다. 제외는 활성 상류에 의해서만 발동(비활성 상류는 흡수 안 함). 결정론.

package graph

// applySupersession applies reins-side Supersedes crisply to the set of active
// counters. For each active upstream counter, it collects the IDs of the
// downstream counters it absorbs into an excluded set, then returns only the
// remaining counters (active and not excluded), preserving registration order.
// Exclusion fires only from an active upstream (an inactive upstream absorbs
// nothing). Deterministic.
func (g *Graph) applySupersession(active []activeCounter) []activeCounter {
	activeIDs := make(map[string]bool, len(active))
	for _, ac := range active {
		activeIDs[ac.node.id] = true
	}
	excluded := make(map[string]bool)
	for _, ac := range active {
		for downID := range g.supersedes[ac.node.id] {
			excluded[downID] = true
		}
	}
	out := make([]activeCounter, 0, len(active))
	for _, ac := range active {
		if excluded[ac.node.id] {
			continue
		}
		out = append(out, ac)
	}
	return out
}
