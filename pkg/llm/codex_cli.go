//ff:type feature=llm type=adapter
//ff:what CodexCLI — `codex exec`(headless 에이전트) 서브프로세스를 호출하는 Backend. HTTP 어댑터와 달리 exec로 돌고, 인증은 codex CLI 자체 로그인(codex login, ChatGPT 구독 / CODEX_HOME)에 위임(REINS 키 없음). codex는 --max-turns가 없어 단일샷을 `-s read-only`(부작용 차단)+no-tools 프리앰블로 *유도*하므로 MaxTurns 필드 없음. 세션 연속(Continue) 모드에서 첫 응답 JSONL의 thread_id를 sid에 운반하므로 *CodexCLI 포인터 리시버. Model 빈 값이면 CLI 기본 모델(-m 생략), Bin 빈 값이면 "codex". claude/grok 백엔드와 동형 — no-tools 프리앰블·runSubprocess seam 공유.

package llm

// CodexCLI is a Backend that shells out to `codex exec` (a headless agent). Unlike
// the HTTP adapters it runs a subprocess; auth is delegated entirely to the codex
// CLI's own login (codex login, ChatGPT subscription / CODEX_HOME — no REINS API-key
// env). codex has no --max-turns, so a single-shot is only *induced* via -s read-only
// (side-effect block) plus the no-tools preamble — hence no MaxTurns field. It is used
// by pointer so the session id (thread_id) survives across Complete calls in Continue
// mode.
type CodexCLI struct {
	Model   string      // "" ⇒ CLI default model (-m omitted); or a model id (e.g. "gpt-5")
	Bin     string      // "" ⇒ "codex"; overridden by env REINS_CODEX_BIN
	Session SessionMode // default Stateless (--ephemeral); Continue is the only opt-in

	sid string // Continue mode carries the first response's thread_id (mutable internal state)
}
