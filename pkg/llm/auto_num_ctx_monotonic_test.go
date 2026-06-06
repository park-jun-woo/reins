//ff:func feature=llm type=helper control=iteration dimension=1
//ff:what TestAutoNumCtxMonotonic — 프롬프트가 길어질수록 컨텍스트 창이 작아지지 않고(단조) 항상 [2048,32768] 범위 안인지 루프로 검증.

package llm

import (
	"strings"
	"testing"
)

// TestAutoNumCtxMonotonic: longer prompts never get a smaller window.
func TestAutoNumCtxMonotonic(t *testing.T) {
	prev := 0
	for n := 0; n <= 200000; n += 5000 {
		got := autoNumCtx(strings.Repeat("z", n))
		if got < prev {
			t.Fatalf("autoNumCtx not monotonic at len=%d: %d < %d", n, got, prev)
		}
		if got < 2048 || got > 32768 {
			t.Fatalf("autoNumCtx(len=%d) = %d out of [2048,32768]", n, got)
		}
		prev = got
	}
}
