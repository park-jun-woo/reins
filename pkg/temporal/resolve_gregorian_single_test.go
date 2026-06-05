//ff:func feature=temporal type=helper control=sequence
//ff:what Absolute+Gregorian 단일 날짜가 그대로 정규 ISO Value로(IsInterval=false, Determined=true) 환원되는지 검증한다.

package temporal

import (
	"testing"
	"time"
)

func TestResolveGregorianSingle(t *testing.T) {
	r := Resolve(Spec{Kind: Absolute, Calendar: Gregorian, Start: "2017-01-10"}, time.Time{})
	if !r.Determined || r.IsInterval || r.Value != "2017-01-10" {
		t.Fatalf("got %+v", r)
	}
}
