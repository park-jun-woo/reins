//ff:func feature=llm type=helper control=iteration dimension=1
//ff:what autoNumCtx — 토큰 수를 ASCII 바이트/4 + 비ASCII rune×2 + reserve(출력 예비분, 0이면 numPredict=2048 폴백)로 근사하고 ollama 컨텍스트 창 계단값(2048→4096→8192→16384→32768)으로 올림 스냅한다. 순수 ASCII는 기존 len/4와 동일, CJK는 rune당 2토큰 보수 추정(len/4 바이트 휴리스틱의 한국어 과소추정 교정 — 단순 rune/4는 더 과소라 금지). 하한 2048·상한 32768 클램프. ★상한 클램프(32768) 탓에 reserve+프롬프트가 32768을 넘으면 컨텍스트 창이 출력을 다 못 담을 수 있다(no silent caps — 호출부 주석 참조).

package llm

// numPredict is the default reserved output token budget when reserve == 0.
const numPredict = 2048

// ctxHi is the largest ollama context-window step (upper clamp).
const ctxHi = 32768

// autoNumCtx estimates the token count of the prompt — ASCII bytes/4 plus 2 tokens
// per non-ASCII rune (CJK runs 1–2.5 tokens/char, so the old bytes/4 heuristic
// underestimated Korean; a plain runes/4 swap would underestimate even more) plus
// the reserve output budget (reserve == 0 ⇒ numPredict default 2048) — and snaps it
// up to the next ollama context-window step (2048 → 4096 → 8192 → 16384 → 32768),
// clamped to [2048, 32768]. A pure-ASCII prompt with the default reserve yields the
// same value as the old len/4. Threading reserve from the caller's effective output
// limit keeps num_ctx large enough to hold a raised output budget (Phase017 A).
func autoNumCtx(prompt string, reserve int) int {
	if reserve == 0 {
		reserve = numPredict
	}
	asciiBytes, nonASCII := 0, 0
	for _, r := range prompt {
		if r < 0x80 {
			asciiBytes++
		} else {
			nonASCII++
		}
	}
	need := asciiBytes/4 + nonASCII*2 + reserve
	for _, step := range []int{2048, 4096, 8192, 16384, 32768} {
		if need <= step {
			return clampCtxLo(step)
		}
	}
	return ctxHi
}
