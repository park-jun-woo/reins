//ff:func feature=temporal type=helper control=sequence
//ff:what Absolute+Gregorian에서 start는 유효하나 end가 파싱 불가하면 Determined=false로 정직 반환하는지 검증한다.

package temporal

import (
	"testing"
	"time"
)

func TestResolveBadEndUndetermined(t *testing.T) {
	r := Resolve(Spec{Kind: Absolute, Calendar: Gregorian, Start: "2017-01-10", End: "nope"}, time.Time{})
	if r.Determined {
		t.Fatalf("got %+v, want undetermined when end is unparseable", r)
	}
}
