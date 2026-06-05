//ff:func feature=quest type=helper control=sequence
//ff:what Progress가 아이템을 상태별로 정확히 집계하는지 검증한다.

package quest

import "testing"

func TestProgressTally(t *testing.T) {
	s := &Session{Items: []*Item{
		{State: TODO},
		{State: TODO},
		{State: PASS},
		{State: DONE},
	}}
	prog := s.Progress()
	if prog[TODO] != 2 {
		t.Errorf("TODO = %d, want 2", prog[TODO])
	}
	if prog[PASS] != 1 || prog[DONE] != 1 {
		t.Errorf("PASS=%d DONE=%d, want 1 each", prog[PASS], prog[DONE])
	}
	if prog[REVIEW] != 0 {
		t.Errorf("REVIEW = %d, want 0", prog[REVIEW])
	}
}
