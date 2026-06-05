//ff:func feature=ground type=helper control=selection
//ff:what HTTPBody — URL의 페이지 본문을 스냅샷에서 읽는다. 첫 읽힘이면 resolver.Fetch를 1회 호출해 본문·에러를 캐시하고, 이후 같은 URL은 캐시만 반환(resolver 추가 호출 0). 에러도 캐시되어 동일 평가 내 결정적으로 재현된다. 호출측은 err≠nil이면 결정론 FAIL Fact로 환원한다.

package ground

// HTTPBody returns the page body at url from the snapshot. On the first read it
// calls resolver.Fetch once and caches the body and any error; subsequent reads of
// the same url return the cache only (zero extra resolver calls). The error is
// cached too, so a failed fetch is reproduced deterministically within the same
// evaluation. The caller reduces a non-nil error to a FAIL Fact.
func (s *Snapshot) HTTPBody(url string) (string, error) {
	if s.bodySeen[url] {
		return s.bodyCache[url], s.bodyErr[url]
	}
	s.bodySeen[url] = true
	body, err := s.resolver.Fetch(url)
	s.bodyCache[url] = body
	s.bodyErr[url] = err
	return body, err
}
