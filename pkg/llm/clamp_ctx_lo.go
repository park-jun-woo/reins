//ff:func feature=llm type=helper control=sequence
//ff:what clampCtxLo — ollama num_ctx 후보값을 하한 2048로 클램프한다(autoNumCtx 스냅 결과의 하한 보장).

package llm

// ctxLo is the smallest ollama context-window step (lower clamp).
const ctxLo = 2048

// clampCtxLo returns step raised to the lower bound ctxLo when it is below it.
func clampCtxLo(step int) int {
	if step < ctxLo {
		return ctxLo
	}
	return step
}
