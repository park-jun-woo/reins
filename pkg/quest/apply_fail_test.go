//ff:func feature=quest type=helper control=iteration dimension=1
//ff:what Apply가 FAIL을 Tries++로 누적하고 MaxTries에서 DONE으로 잠그는지(소진 종료) 검증한다.

package quest

import "testing"

func TestApplyFailRatchetToDone(t *testing.T) {
	it := &Item{Key: "x", State: TODO}
	for i := 0; i < MaxTries-1; i++ {
		Apply(it, Verdict{Outcome: OutFail}, "now")
		if it.State != TODO {
			t.Fatalf("after %d fails state = %s, want TODO", i+1, it.State)
		}
	}
	Apply(it, Verdict{Outcome: OutFail}, "now")
	if it.State != DONE {
		t.Fatalf("after MaxTries state = %s, want DONE", it.State)
	}
	if it.Tries != MaxTries {
		t.Fatalf("tries = %d, want %d", it.Tries, MaxTries)
	}
}
