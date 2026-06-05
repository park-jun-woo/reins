package quest

import "testing"

func TestApplyPassLocks(t *testing.T) {
	it := &Item{Key: "x", State: TODO}
	Apply(it, Verdict{Outcome: OutPass}, "now")
	if it.State != PASS || !it.State.Terminal() {
		t.Fatalf("state = %s, want locked PASS", it.State)
	}
}

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
