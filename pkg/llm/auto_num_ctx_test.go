//ff:func feature=llm type=helper control=iteration dimension=1
//ff:what TestAutoNumCtx — 프롬프트 바이트 길이→컨텍스트 창 계단값 스냅과 [2048,32768] 클램프를 테이블로 검증(need = len/4 + numPredict).

package llm

import (
	"strings"
	"testing"
)

// TestAutoNumCtx checks the byte-length → context-window step snapping plus the
// [2048, 32768] clamp. need = len/4 + numPredict(2048).
func TestAutoNumCtx(t *testing.T) {
	cases := []struct {
		name      string
		promptLen int
		want      int
	}{
		// need <= 2048 ⇒ clamp to lower bound 2048. len 0 ⇒ need 2048.
		{"empty", 0, 2048},
		// boundary: need exactly 4096 ⇒ stays 4096. len/4 = 2048 ⇒ len 8192.
		{"exactly-4096", (4096 - numPredict) * 4, 4096},
		// boundary: need exactly 8192 ⇒ 8192. len/4 = 6144 ⇒ len 24576.
		{"exactly-8192", (8192 - numPredict) * 4, 8192},
		// boundary: need exactly 16384 ⇒ 16384.
		{"exactly-16384", (16384 - numPredict) * 4, 16384},
		// boundary: need exactly 32768 ⇒ 32768.
		{"exactly-32768", (32768 - numPredict) * 4, 32768},
		// over the top step ⇒ clamp to hi 32768.
		{"over-hi", 1 << 20, 32768},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := autoNumCtx(strings.Repeat("a", c.promptLen))
			if got != c.want {
				t.Fatalf("autoNumCtx(len=%d) = %d, want %d", c.promptLen, got, c.want)
			}
		})
	}
}
