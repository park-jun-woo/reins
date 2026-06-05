//ff:func feature=temporal type=helper control=sequence
//ff:what Relative 명세가 ref+OffsetDays로 정규화되는지 검증한다(어제 = ref-1일, Determined=true).

package temporal

import (
	"testing"
	"time"
)

func TestResolveRelative(t *testing.T) {
	ref := time.Date(2026, 6, 5, 0, 0, 0, 0, time.UTC)
	r := Resolve(Spec{Kind: Relative, OffsetDays: -1}, ref)
	if !r.Determined || r.Value != "2026-06-04" {
		t.Fatalf("got %+v, want 2026-06-04 determined", r)
	}
}
