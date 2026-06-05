//ff:func feature=graph type=helper control=selection
//ff:what selectRootCause — 잔존 카운터에서 결정타(root cause)와 곁가지 Fail들을 가른다. 결정타 = 잔존 Fail 중 supersession 위상 최상류(supersedesReach 최대), 동률·무관계면 등록순 첫째. 나머지 Fail은 siblings로 병기. 잔존 Fail이 없으면(REVIEW) 첫 잔존 카운터를 결정타로 둔다. remaining은 등록순이라 결정론.

package graph

import "github.com/park-jun-woo/reins/pkg/gate"

// selectRootCause partitions the remaining counters into the root cause (decisive
// blow) and the other Fail counters (siblings). The root cause is the remaining
// Fail with the greatest supersession reach (topologically uppermost); ties or
// no-relation fall back to registration order (first). Remaining Fail counters
// other than the root cause are returned as siblings. When no remaining Fail exists
// (a REVIEW), the first remaining counter is used as the root cause. remaining is in
// registration order, so the choice is deterministic.
func (g *Graph) selectRootCause(remaining []activeCounter) (root activeCounter, siblings []activeCounter) {
	fails := make([]activeCounter, 0, len(remaining))
	for _, ac := range remaining {
		if ac.node.level == gate.LevelFail {
			fails = append(fails, ac)
		}
	}
	if len(fails) == 0 {
		// REVIEW: no Fail; the first remaining counter is the focal point.
		if len(remaining) > 0 {
			return remaining[0], nil
		}
		return activeCounter{}, nil
	}
	rootIdx := 0
	bestReach := g.supersedesReach(fails[0].node.id)
	for i := 1; i < len(fails); i++ {
		if r := g.supersedesReach(fails[i].node.id); r > bestReach {
			bestReach, rootIdx = r, i
		}
	}
	root = fails[rootIdx]
	for i, ac := range fails {
		if i != rootIdx {
			siblings = append(siblings, ac)
		}
	}
	return root, siblings
}
