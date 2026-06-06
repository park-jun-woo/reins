//ff:func feature=llm type=helper control=iteration dimension=1
//ff:what autoNumCtx — 프롬프트 바이트 길이로 토큰 수를 근사(≈len/4 + num_predict 여유분)하고 ollama 컨텍스트 창 계단값(2048→4096→8192→16384→32768)으로 올림 스냅한다. 하한 2048·상한 32768 클램프. 긴 기사면 큰 창, 짧으면 작은 창으로 메모리·속도 균형.

package llm

// numPredict is the reserved output token budget added to the input estimate.
const numPredict = 2048

// ctxHi is the largest ollama context-window step (upper clamp).
const ctxHi = 32768

// autoNumCtx estimates the token count of the prompt (≈ len/4 bytes-per-token plus
// the num_predict output reserve) and snaps it up to the next ollama context-window
// step (2048 → 4096 → 8192 → 16384 → 32768), clamped to [2048, 32768].
func autoNumCtx(prompt string) int {
	need := len(prompt)/4 + numPredict
	for _, step := range []int{2048, 4096, 8192, 16384, 32768} {
		if need <= step {
			return clampCtxLo(step)
		}
	}
	return ctxHi
}
