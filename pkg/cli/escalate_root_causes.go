//ff:func feature=cli type=helper control=iteration dimension=1
//ff:what escalateRootCauses — EscalateOn 슬라이스를 O(1) 조회용 set으로 변환. nil·빈 슬라이스는 빈 set을 내므로 에스컬레이션이 절대 발화하지 않는다.

package cli

// escalateRootCauses builds a set from the EscalateOn slice for O(1) lookup. A nil
// or empty slice yields an empty set, so escalation never fires.
func escalateRootCauses(ids []string) map[string]bool {
	m := make(map[string]bool, len(ids))
	for _, id := range ids {
		m[id] = true
	}
	return m
}
