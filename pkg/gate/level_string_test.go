//ff:func feature=gate type=helper control=sequence
//ff:what Level.String()이 LevelFail→"FAIL", LevelReview→"REVIEW"를 내는지 검증한다.

package gate

import "testing"

func TestLevel_String(t *testing.T) {
	if got := LevelFail.String(); got != "FAIL" {
		t.Errorf("LevelFail.String() = %q, want FAIL", got)
	}
	if got := LevelReview.String(); got != "REVIEW" {
		t.Errorf("LevelReview.String() = %q, want REVIEW", got)
	}
}
