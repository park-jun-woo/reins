//ff:func feature=ground type=helper control=selection
//ff:what MXResolves — 도메인이 수신 가능한지(MX 레코드 존재)를 스냅샷에서 읽는다. 첫 읽힘이면 resolver.LookupMX를 1회 호출해 bool·에러를 캐시하고, 이후 같은 도메인은 캐시만 반환(resolver 추가 호출 0). 에러도 캐시되어 동일 평가 내 결정적. 호출측은 err≠nil이면 결정론 FAIL Fact로 환원한다.

package ground

// MXResolves reports whether domain is deliverable (has an MX record) from the
// snapshot. On the first read it calls resolver.LookupMX once and caches the bool
// and any error; subsequent reads of the same domain return the cache only (zero
// extra resolver calls). The error is cached too, so a failed lookup is reproduced
// deterministically within the same evaluation. The caller reduces a non-nil error
// to a FAIL Fact.
func (s *Snapshot) MXResolves(domain string) (bool, error) {
	if s.mxSeen[domain] {
		return s.mxCache[domain], s.mxErr[domain]
	}
	s.mxSeen[domain] = true
	ok, err := s.resolver.LookupMX(domain)
	s.mxCache[domain] = ok
	s.mxErr[domain] = err
	return ok, err
}
