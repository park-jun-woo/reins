//ff:func feature=llm type=helper control=iteration dimension=1
//ff:what TestAutoNumCtxCJK — 같은 글자 수의 한국어 프롬프트가 영문 이상(≥)의 컨텍스트 창을 받는지(비ASCII rune×2 보수 추정), 과소추정이 실제로 교정되는 구간에서는 엄격히 큰지 검증한다.

package llm

import (
	"strings"
	"testing"
)

// TestAutoNumCtxCJK: a Korean prompt of the same character count is never sized
// below the English one (2 tokens per non-ASCII rune), and in the range where the
// old bytes/4 heuristic underestimated, the Korean window is strictly larger.
func TestAutoNumCtxCJK(t *testing.T) {
	for _, n := range []int{0, 100, 500, 2000, 8000, 30000} {
		ko := autoNumCtx(strings.Repeat("가", n))
		en := autoNumCtx(strings.Repeat("a", n))
		if ko < en {
			t.Fatalf("chars=%d: korean ctx %d < english ctx %d", n, ko, en)
		}
	}
	// 2000 chars: korean need = 2000*2+2048 = 6048 → 8192; english need = 2548 → 4096.
	if ko, en := autoNumCtx(strings.Repeat("가", 2000)), autoNumCtx(strings.Repeat("a", 2000)); ko <= en {
		t.Fatalf("chars=2000: korean ctx %d not strictly greater than english %d", ko, en)
	}
}
