//ff:func feature=temporal type=helper control=sequence
//ff:what parseGregorian이 파싱 불가 성분엔 ok=false를 내는지 검증한다.

package temporal

import "testing"

func TestParseGregorianInvalid(t *testing.T) {
	if _, ok := parseGregorian("not-a-date"); ok {
		t.Fatal("unparseable component should return ok=false")
	}
}
