//ff:func feature=graph type=helper control=iteration dimension=1
//ff:what 테스트 헬퍼 — Verdict의 Facts에서 Where 값을 카운트 맵으로 뽑는다. 동치 테스트가 Fact 집합을 순서무시로 비교하는 입력.

package graph

import "github.com/park-jun-woo/reins/pkg/quest"

// factWheres extracts a count map of the Where values from a Verdict's Facts.
func factWheres(v quest.Verdict) map[string]int {
	out := make(map[string]int)
	for _, f := range v.Facts {
		out[f.Where]++
	}
	return out
}
