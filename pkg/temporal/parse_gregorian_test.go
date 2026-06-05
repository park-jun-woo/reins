//ff:func feature=temporal type=helper control=sequence
//ff:what parseGregorian이 유효 ISO 날짜를 정규 형태·ok=true로 돌려주는지 검증한다.

package temporal

import "testing"

func TestParseGregorianValid(t *testing.T) {
	got, ok := parseGregorian("2017-01-10")
	if !ok {
		t.Fatal("valid ISO date should parse")
	}
	if got != "2017-01-10" {
		t.Fatalf("got %q, want canonical 2017-01-10", got)
	}
}
