//ff:func feature=llm type=adapter control=sequence
//ff:what newCodexCLI — model 토큰 + REINS_CODEX_* env로 *CodexCLI를 짓는다. model=="default"→Model:""(CLI 기본 모델 sentinel; FromFlag 빈-model 거부 우회). env REINS_CODEX_BIN→Bin. REINS_CODEX_SESSION=="continue"→Continue 옵트인; 그 외(미설정·"stateless"·인식 불가)→Stateless 안전 폴백(인식 불가 값을 조용히 외부 세션 ID로 해석하지 않음 — "continue" 오타가 가짜 세션을 못 만듦). newGrokCLI/newClaudeCLI와 동형.

package llm

import "os"

// newCodexCLI builds a *CodexCLI from the model token plus REINS_CODEX_* env.
//
//	model == "default" ⇒ Model:"" (use the CLI default model; sentinel that bypasses
//	                      FromFlag's empty-model rejection)
//	env REINS_CODEX_BIN     ⇒ Bin
//	env REINS_CODEX_SESSION ⇒ "continue" → Continue (session-continuity opt-in)
//	                          "" / "stateless" / unrecognized → Stateless (safe default;
//	                          an unknown value is never silently read as an external
//	                          session id, so a "continue" typo cannot forge a session)
func newCodexCLI(model string) *CodexCLI {
	if model == "default" {
		model = ""
	}
	c := &CodexCLI{
		Model: model,
		Bin:   os.Getenv("REINS_CODEX_BIN"),
	}
	if os.Getenv("REINS_CODEX_SESSION") == "continue" {
		c.Session = SessionMode{Kind: Continue}
	}
	return c
}
