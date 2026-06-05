//ff:func feature=temporal type=helper control=sequence
//ff:what Absolute+Gregorian 기간(start+end)이 "start/end" Value·IsInterval=true·Determined=true로 환원되는지 검증한다.

package temporal

import (
	"testing"
	"time"
)

func TestResolveGregorianInterval(t *testing.T) {
	r := Resolve(Spec{Kind: Absolute, Calendar: Gregorian, Start: "2017-01-10", End: "2017-01-14"}, time.Time{})
	if !r.Determined || !r.IsInterval || r.Value != "2017-01-10/2017-01-14" {
		t.Fatalf("got %+v", r)
	}
}
