//ff:func feature=graph type=helper control=iteration dimension=1
//ff:what 테스트 헬퍼 — 두 Where-카운트 맵이 순서무시로 동일한지 판정한다. 동치 테스트의 Fact 집합 일치 검사에 쓴다.

package graph

// equalStringSets reports whether two Where-count maps are equal (order-free).
func equalStringSets(a, b map[string]int) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}
