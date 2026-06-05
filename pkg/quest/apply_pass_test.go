//ff:func feature=quest type=helper control=sequence
//ff:what Apply가 PASS verdict로 아이템을 terminal PASS로 잠그는지 검증한다.

package quest

import "testing"

func TestApplyPassLocks(t *testing.T) {
	it := &Item{Key: "x", State: TODO}
	Apply(it, Verdict{Outcome: OutPass}, "now")
	if it.State != PASS || !it.State.Terminal() {
		t.Fatalf("state = %s, want locked PASS", it.State)
	}
	if it.CollectedAt != "now" {
		t.Errorf("CollectedAt = %q, want %q", it.CollectedAt, "now")
	}
}
