//ff:func feature=cli type=helper control=sequence
//ff:what composeSystem — 전역 system 프롬프트에 규칙별 코칭을 결합한다. 빈 전역은 reins 기본 프롬프트(fallbackSystem)로, 빈 코칭은 결합 없이 전역만 반환한다.

package cli

// fallbackSystem is the generic system prompt used when LoopOptions.System is empty.
const fallbackSystem = "You produce a submission a deterministic gate will judge. " +
	"Output only the payload in the exact format the prompt specifies; no prose."

// composeSystem combines the global system prompt with optional rule-specific
// coaching. An empty global falls back to the generic reins prompt; empty coaching
// adds nothing.
func composeSystem(global, ruleCoach string) string {
	if global == "" {
		global = fallbackSystem
	}
	if ruleCoach == "" {
		return global
	}
	return global + "\n" + ruleCoach
}
