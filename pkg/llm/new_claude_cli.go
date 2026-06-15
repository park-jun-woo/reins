//ff:func feature=llm type=adapter control=sequence
//ff:what newClaudeCLI — model 토큰 + REINS_CLAUDE_* env로 *ClaudeCLI를 짓는다. model=="default"→Model:""(CLI 기본 모델; FromFlag 빈-model 거부 우회 sentinel). env REINS_CLAUDE_BIN→Bin. REINS_CLAUDE_SESSION=="continue"→Continue 옵트인; 그 외(미설정·"stateless"·인식 불가)→Stateless 안전 폴백(인식 불가 값을 조용히 외부 세션 UUID로 해석하지 않음 — "continue" 오타가 가짜 세션을 못 만듦).

package llm

import "os"

// newClaudeCLI builds a *ClaudeCLI from the model token plus REINS_CLAUDE_* env.
//
//	model == "default" ⇒ Model:"" (use the CLI default model; sentinel that bypasses
//	                      FromFlag's empty-model rejection)
//	env REINS_CLAUDE_BIN     ⇒ Bin
//	env REINS_CLAUDE_SESSION ⇒ "continue" → Continue (session-continuity opt-in)
//	                           "" / "stateless" / unrecognized → Stateless (safe default;
//	                           an unknown value is never silently read as an external
//	                           session UUID, so a "continue" typo cannot forge a session)
func newClaudeCLI(model string) *ClaudeCLI {
	if model == "default" {
		model = ""
	}
	c := &ClaudeCLI{
		Model: model,
		Bin:   os.Getenv("REINS_CLAUDE_BIN"),
	}
	if os.Getenv("REINS_CLAUDE_SESSION") == "continue" {
		c.Session = SessionMode{Kind: Continue}
	}
	return c
}
