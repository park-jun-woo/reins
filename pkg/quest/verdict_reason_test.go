//ff:func feature=quest type=helper control=sequence
//ff:what Verdict.Reason이 Facts가 없으면 Outcome만 렌더하는지 검증한다.

package quest

import "testing"

func TestReasonNoFacts(t *testing.T) {
	v := Verdict{Outcome: OutPass}
	if got := v.Reason(); got != "PASS" {
		t.Errorf("Reason() = %q, want %q", got, "PASS")
	}
}
