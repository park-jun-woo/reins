//ff:type feature=cli type=model
//ff:what LoopOptions — in-process generate→gate→retry 에이전트 루프 설정. LLM은 생성자(L0)만 — PASS 잠금 권한은 게이트에. DefaultModel(--model 기본)·System(전역 system 프롬프트)·RuleSystem(rule ID→FAIL 시 추가 코칭)·LLM(주입 backend, 비-nil이면 --model 무시).

package cli

import (
	"github.com/park-jun-woo/reins/pkg/llm"
)

// LoopOptions configures the in-process generate→gate→retry loop. The LLM is
// the generator (L0) only — PASS lock authority stays with the gate.
type LoopOptions struct {
	// DefaultModel is the --model default. Empty ⇒ "ollama:gemma4:e4b".
	DefaultModel string
	// System is the global system prompt.
	System string
	// RuleSystem maps a toulmin rule ID (verdict.RootCause) to extra system
	// guidance appended when the previous attempt FAILed on that rule.
	RuleSystem map[string]string
	// LLM, when non-nil, is used as the backend and --model is ignored (for tests
	// or a fixed backend).
	LLM llm.Backend
	// Escalate, when non-nil, is a stronger fallback backend. Once an item's FAIL
	// carries a RootCause listed in EscalateOn (a capability-bound signal — a
	// semantic mismatch, not a format slip), the item is retried with Escalate for
	// its remaining tries (latched on for that item). The gate still holds sole PASS
	// authority — escalation only changes which generator (L0) is asked. nil ⇒ no
	// escalation (backward compatible). Cost note: a slow/paid backend here is
	// invoked only on the residual the primary cannot crack.
	Escalate llm.Backend
	// EscalateOn is the set of FAIL RootCause IDs that promote an item to Escalate.
	// Empty ⇒ never escalate (even when Escalate is set), so format/shape failures
	// the consumer leaves out stay on the cheap primary.
	EscalateOn []string
}
