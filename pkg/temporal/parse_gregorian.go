//ff:func feature=temporal type=helper control=sequence
//ff:what 그레고리력 성분 한 개("2006-01-02")를 검증해 정규 ISO 문자열과 성공여부를 돌려준다. 파싱 실패 시 ok=false(호출측이 Determined=false로 매핑). 표준 time만 사용.

package temporal

import "time"

// parseGregorian validates a single Gregorian component as an ISO calendar date
// ("2006-01-02") and returns its canonical re-formatted form. It returns ok=false
// for any unparseable component so Resolve can map it to Determined=false. v1 uses
// the standard library only; non-Gregorian calendar conversion is v2 (deferred).
func parseGregorian(s string) (string, bool) {
	t, err := time.Parse(isoLayout, s)
	if err != nil {
		return "", false
	}
	return t.Format(isoLayout), true
}
