//ff:type feature=llm type=adapter
//ff:what GeminiCLI — `gemini -p`(headless 비대화형) 서브프로세스를 호출하는 Backend. HTTP 어댑터(Gemini, gemini.go)와 달리 서브프로세스로 돌고, 인증은 gemini CLI 자체 로그인(Google 계정 OAuth / 그 프로세스의 GEMINI_API_KEY·Vertex)에 위임(REINS 키 없음 — HTTP gemini:만 GEMINI_API_KEY를 읽는다). gemini는 --max-turns가 없어 단일샷을 `--approval-mode plan`(read-only·부작용 차단)+no-tools 프리앰블로 *유도*하므로 MaxTurns 필드 없음. Continue 모드에서 reins가 발급한 세션 UUID를 sid에 운반(1차 --session-id, 2차+ --resume latest)하므로 *GeminiCLI 포인터 리시버. Model 빈 값이면 CLI 기본 모델(-m 생략), Bin 빈 값이면 "gemini". claude/grok/codex 백엔드와 동형 — no-tools 프리앰블·runSubprocess seam 공유.

package llm

// GeminiCLI is a Backend that shells out to `gemini -p` (a headless, non-interactive
// run). Unlike the HTTP adapter (Gemini, gemini.go) it runs a subprocess; auth is
// delegated entirely to the gemini CLI's own login (Google OAuth, or that process's
// GEMINI_API_KEY / Vertex — no REINS key; only the HTTP gemini: backend reads
// GEMINI_API_KEY). gemini has no --max-turns, so a single-shot is only *induced* via
// --approval-mode plan (read-only, side-effect block) plus the no-tools preamble —
// hence no MaxTurns field. It is used by pointer so the reins-issued session UUID
// survives across Complete calls in Continue mode (--session-id then --resume latest).
type GeminiCLI struct {
	Model   string      // "" ⇒ CLI default model (-m omitted); or a model id (e.g. "gemini-2.5-pro")
	Bin     string      // "" ⇒ "gemini"; overridden by env REINS_GEMINI_BIN
	Session SessionMode // default Stateless; Continue is the only opt-in

	sid string // Continue mode carries the reins-issued session UUID (mutable internal state)
}
