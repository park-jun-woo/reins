//ff:type feature=llm type=model
//ff:what sessionKind — ClaudeCLI 세션 동작 선택자. Stateless(기본, 매 Complete 독립 — HTTP 백엔드와 동의)·Continue(첫 응답 session_id를 다음 호출 --resume로 운반, loop 한 실행 동안 대화 누적). 두 값뿐 — 외부 UUID 핀 같은 추가 모드는 범위 밖.

package llm

// sessionKind selects whether each Complete is independent (Stateless) or carries
// the claude session forward across calls (Continue).
type sessionKind int

const (
	// Stateless makes every Complete fully independent (--no-session-persistence),
	// matching the HTTP backends. Default.
	Stateless sessionKind = iota
	// Continue carries the first response's session_id into later calls (--resume),
	// accumulating the claude conversation across a loop run (opt-in escape hatch).
	Continue
)
