//ff:func feature=ground type=helper control=selection
//ff:what defaultResolver.LookupMX — 실제 net.LookupMX로 도메인의 MX 레코드 존재를 본다. 레코드가 하나라도 있으면 (true, nil), 빈 레코드 집합은 (false, nil), 조회 실패는 (false, err). 테스트는 fake를 주입해 네트워크를 피한다.

package ground

import "net"

// LookupMX reports whether domain has at least one MX record. An empty record set is
// (false, nil); a lookup failure is (false, err).
func (r defaultResolver) LookupMX(domain string) (bool, error) {
	recs, err := net.LookupMX(domain)
	if err != nil {
		return false, err
	}
	return len(recs) > 0, nil
}
