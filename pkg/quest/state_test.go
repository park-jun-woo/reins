//ff:func feature=quest type=helper control=iteration dimension=1
//ff:what State.Terminal이 TODO만 비종단이고 나머지 5개는 잠금(종단)으로 보고하는지 검증한다(래칫 단방향).

package quest

import "testing"

func TestTerminal(t *testing.T) {
	cases := map[State]bool{
		TODO:    false,
		PASS:    true,
		REVIEW:  true,
		DONE:    true,
		SKIPPED: true,
		BLOCKED: true,
	}
	for s, want := range cases {
		if got := s.Terminal(); got != want {
			t.Errorf("%s.Terminal() = %v, want %v", s, got, want)
		}
	}
}
