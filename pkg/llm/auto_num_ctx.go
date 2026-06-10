//ff:func feature=llm type=helper control=iteration dimension=1
//ff:what autoNumCtx — 토큰 수를 ASCII 바이트/4 + 비ASCII rune×2 + num_predict 여유분으로 근사하고 ollama 컨텍스트 창 계단값(2048→4096→8192→16384→32768)으로 올림 스냅한다. 순수 ASCII는 기존 len/4와 동일, CJK는 rune당 2토큰 보수 추정(len/4 바이트 휴리스틱의 한국어 과소추정 교정 — 단순 rune/4는 더 과소라 금지). 하한 2048·상한 32768 클램프.

package llm

// numPredict is the reserved output token budget added to the input estimate.
const numPredict = 2048

// ctxHi is the largest ollama context-window step (upper clamp).
const ctxHi = 32768

// autoNumCtx estimates the token count of the prompt — ASCII bytes/4 plus 2 tokens
// per non-ASCII rune (CJK runs 1–2.5 tokens/char, so the old bytes/4 heuristic
// underestimated Korean; a plain runes/4 swap would underestimate even more) plus
// the num_predict output reserve — and snaps it up to the next ollama
// context-window step (2048 → 4096 → 8192 → 16384 → 32768), clamped to
// [2048, 32768]. A pure-ASCII prompt yields the same value as the old len/4.
func autoNumCtx(prompt string) int {
	asciiBytes, nonASCII := 0, 0
	for _, r := range prompt {
		if r < 0x80 {
			asciiBytes++
		} else {
			nonASCII++
		}
	}
	need := asciiBytes/4 + nonASCII*2 + numPredict
	for _, step := range []int{2048, 4096, 8192, 16384, 32768} {
		if need <= step {
			return clampCtxLo(step)
		}
	}
	return ctxHi
}
