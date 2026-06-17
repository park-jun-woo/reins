//ff:func feature=llm type=helper control=sequence
//ff:what TestAutoNumCtxJustOverLowerBound — need가 2048을 막 넘는 프롬프트는 하한 클램프가 아니라 다음 계단값(4096)으로 스냅되는지 검증.

package llm

import (
	"strings"
	"testing"
)

// TestAutoNumCtxJustOverLowerBound: a prompt whose need barely exceeds 2048 snaps
// up to the next step (4096), not the lower clamp.
func TestAutoNumCtxJustOverLowerBound(t *testing.T) {
	// len 8 ⇒ len/4 = 2 ⇒ need 2050 > 2048 ⇒ next step 4096.
	if got := autoNumCtx(strings.Repeat("x", 8), 0); got != 4096 {
		t.Fatalf("autoNumCtx(len=8) = %d, want 4096", got)
	}
}
