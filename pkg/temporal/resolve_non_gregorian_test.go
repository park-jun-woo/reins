//ff:func feature=temporal type=helper control=sequence
//ff:what 비그레고리 역법(Persian)은 v1에서 변환하지 않고 Determined=false로 정직 반환하는지 검증한다.

package temporal

import (
	"testing"
	"time"
)

func TestResolveNonGregorianUndetermined(t *testing.T) {
	if Resolve(Spec{Kind: Absolute, Calendar: Persian, Start: "1395-10-21"}, time.Time{}).Determined {
		t.Fatal("non-gregorian should be undetermined")
	}
}
