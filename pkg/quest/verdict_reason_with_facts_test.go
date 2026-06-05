//ff:func feature=quest type=helper control=sequence
//ff:what Verdict.Reason이 Facts가 있으면 rule/where/expected/got를 한 줄로 렌더하는지 검증한다.

package quest

import "testing"

func TestReasonWithFacts(t *testing.T) {
	v := Verdict{
		Outcome: OutFail,
		Facts: []Fact{
			{Rule: "R1", Where: "title", Expected: "x", Actual: "y"},
			{Where: "body"},
		},
	}
	want := "R1: title (expected x, got y); body"
	if got := v.Reason(); got != want {
		t.Errorf("Reason() = %q, want %q", got, want)
	}
}
