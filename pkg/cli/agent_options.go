//ff:type feature=cli type=model
//ff:what AgentOptions — in-process generate→gate→retry 에이전트 루프 설정. LLM은 생성자(L0)만 — PASS 잠금 권한은 게이트에. DefaultModel(--model 기본)·System(전역 system 프롬프트)·RuleSystem(rule ID→FAIL 시 추가 코칭)·LLM(주입 backend, 비-nil이면 --model 무시).

package cli

import (
	"github.com/park-jun-woo/reins/pkg/llm"
)

// AgentOptions configures the in-process generate→gate→retry agent loop. The LLM is
// the generator (L0) only — PASS lock authority stays with the gate.
type AgentOptions struct {
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
}
