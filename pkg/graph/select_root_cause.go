//ff:func feature=graph type=helper control=selection
//ff:what selectRootCause — 잔존 카운터에서 결정타(root cause)와 곁가지를 가른다. 결정타 = 잔존 Fail 중 supersession 위상 최상류(supersedesReach 최대), 동률·무관계면 등록순 첫째. 잔존 Fail이 없으면(REVIEW) 첫 잔존 카운터를 결정타로 둔다. 결정타를 제외한 나머지 잔존 전부(레벨 무관 — FAIL root와 함께 잔존한 REVIEW 카운터도 포함, 피드백 패리티)가 siblings. remaining은 등록순이라 결정론.

package graph

import "github.com/park-jun-woo/reins/pkg/gate"

// selectRootCause partitions the remaining counters into the root cause (decisive
// blow) and the other remaining counters (siblings). The root cause is the
// remaining Fail with the greatest supersession reach (topologically uppermost);
// ties or no-relation fall back to registration order (first). When no remaining
// Fail exists (a REVIEW), the first remaining counter is used as the root cause.
// Every remaining counter except the root — regardless of level, so a REVIEW
// counter surviving next to a FAIL root is listed too (feedback parity) — is
// returned as a sibling. remaining is in registration order, so the choice is
// deterministic.
func (g *Graph) selectRootCause(remaining []activeCounter) (root activeCounter, siblings []activeCounter) {
	if len(remaining) == 0 {
		return activeCounter{}, nil
	}
	rootPos, bestReach := -1, -1
	for i, ac := range remaining {
		if ac.node.level != gate.LevelFail {
			continue
		}
		if r := g.supersedesReach(ac.node.id); rootPos == -1 || r > bestReach {
			rootPos, bestReach = i, r
		}
	}
	if rootPos == -1 {
		// REVIEW: no Fail; the first remaining counter is the focal point.
		rootPos = 0
	}
	for i, ac := range remaining {
		if i != rootPos {
			siblings = append(siblings, ac)
		}
	}
	return remaining[rootPos], siblings
}
