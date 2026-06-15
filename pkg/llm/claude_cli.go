//ff:type feature=llm type=adapter
//ff:what ClaudeCLI — `claude -p`(headless print) 서브프로세스를 호출하는 Backend. HTTP 어댑터와 달리 exec로 돈다; 인증은 claude CLI 자체 로그인에 위임(REINS API 키 없음). 세션 연속(Continue) 모드에서 첫 응답의 session_id를 sid에 운반하므로 *ClaudeCLI 포인터 리시버. Model 빈 값이면 CLI 기본 모델, Bin 빈 값이면 "claude", MaxTurns 0이면 1(단일샷 L0 강제).

package llm

// ClaudeCLI is a Backend that shells out to `claude -p` (headless print mode).
// Unlike the HTTP adapters it runs a subprocess; auth is delegated entirely to
// the claude CLI's own login (no REINS API-key env). It is used by pointer so the
// session id survives across Complete calls in Continue mode.
type ClaudeCLI struct {
	Model    string      // "" ⇒ CLI default model (--model omitted); alias (sonnet/opus/haiku/fable) or full id
	Bin      string      // "" ⇒ "claude"; overridden by env REINS_CLAUDE_BIN
	Session  SessionMode // default Stateless; Continue is the only opt-in
	MaxTurns int         // 0 ⇒ 1 (single-shot L0 generation, blocks the agentic tool loop)

	sid string // Continue mode carries the first response's session_id (mutable internal state)
}
