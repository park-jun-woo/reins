//ff:func feature=quest type=helper control=sequence
//ff:what lock이 아이템을 주어진 terminal 상태로 잠그고 verdict·reason·수집 시각을 기록하는지 검증한다.

package quest

import "testing"

func TestLockRecords(t *testing.T) {
	it := &Item{Key: "x", State: TODO}
	v := Verdict{Outcome: OutReview, Facts: []Fact{{Rule: "r", Where: "body", Expected: "a", Actual: "b"}}}
	lock(it, REVIEW, v, "2026-06-05T00:00:00Z")
	if it.State != REVIEW {
		t.Errorf("State = %s, want REVIEW", it.State)
	}
	if it.Verdict != string(OutReview) {
		t.Errorf("Verdict = %q, want %q", it.Verdict, OutReview)
	}
	if it.Reason != v.Reason() {
		t.Errorf("Reason = %q, want %q", it.Reason, v.Reason())
	}
	if it.CollectedAt != "2026-06-05T00:00:00Z" {
		t.Errorf("CollectedAt = %q, want injected now", it.CollectedAt)
	}
}
