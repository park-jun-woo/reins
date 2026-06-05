//ff:func feature=temporal type=helper control=sequence
//ff:what 파싱 불가한 그레고리 성분("nope")이 Determined=false로 반환되는지 검증한다.

package temporal

import (
	"testing"
	"time"
)

func TestResolveUnparseableUndetermined(t *testing.T) {
	if Resolve(Spec{Kind: Absolute, Calendar: Gregorian, Start: "nope"}, time.Time{}).Determined {
		t.Fatal("unparseable should be undetermined")
	}
}
