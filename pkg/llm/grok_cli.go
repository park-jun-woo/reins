//ff:type feature=llm type=adapter
//ff:what GrokCLI — `grok -p`(headless single-turn) 서브프로세스를 호출하는 Backend. HTTP 어댑터(xai)와 달리 exec로 돌고, 인증은 grok CLI 자체 로그인(grok login, 구독·OAuth)에 위임(REINS API 키 없음). 세션 연속(Continue) 모드에서 첫 응답의 sessionId를 sid에 운반하므로 *GrokCLI 포인터 리시버. Model 빈 값이면 CLI 기본 모델(-m 생략), Bin 빈 값이면 "grok", MaxTurns 0이면 1(단일샷 L0 강제). claude 백엔드와 동형 — no-tools 프리앰블·runSubprocess seam 공유.

package llm

// GrokCLI is a Backend that shells out to `grok -p` (headless single-turn). Unlike
// the HTTP adapters it runs a subprocess; auth is delegated entirely to the grok
// CLI's own login (grok login, subscription/OAuth — no REINS API-key env). It is
// used by pointer so the session id survives across Complete calls in Continue mode.
type GrokCLI struct {
	Model    string      // "" ⇒ CLI default model (-m omitted); or a model id (e.g. "grok-4")
	Bin      string      // "" ⇒ "grok"; overridden by env REINS_GROK_BIN
	Session  SessionMode // default Stateless; Continue is the only opt-in
	MaxTurns int         // 0 ⇒ 1 (single-shot L0 generation, blocks the agentic tool loop)

	sid string // Continue mode carries the first response's sessionId (mutable internal state)
}
