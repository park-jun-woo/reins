//ff:func feature=cli type=helper control=sequence level=error
//ff:what generatePayload — runLoopItem 한 시도의 L0 생성 단계. composeSystem(System,ruleCoach)으로 system을 짜고 def.Render→backend.Complete로 페이로드를 만든다. Complete 에러는 Phase012 강등: backendErrorVerdict로 합성→applyVerdict로 래칫(Tries++, MaxTries에서 DONE 잠금)→renderVerdict 출력 후 handled=true로 반환해 호출부가 루프를 continue하게 한다(루프 abort 아님). def.Render·applyVerdict(영속화) 에러만 err로 전파해 fail-fast. 인프라 에러는 내용 비평이 아니므로 피드백을 되먹이지 않는다(호출부가 피드백 비갱신).

package cli

import (
	"io"

	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/llm"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// generatePayload runs one attempt's L0 generation: it composes the system prompt,
// renders the item, and calls backend.Complete. A def.Render error propagates (err).
// A backend.Complete error is demoted (Phase012): a synthetic backend-error verdict
// is ratcheted (Tries++, locking DONE at MaxTries) and rendered, then handled=true is
// returned so the caller continues the loop — generation failure is a retryable item
// FAIL, not a run abort. A persistence failure inside applyVerdict is still fatal
// (err). No feedback is fed back for an infra error.
func generatePayload(def gate.Definition, opts *LoopOptions, backend llm.Backend, ruleCoach, feedback string, s *quest.Session, it *quest.Item, outPath, sessionPath string, out io.Writer) (raw string, handled bool, err error) {
	system := composeSystem(opts.System, ruleCoach)
	prompt, err := def.Render(s, it)
	if err != nil {
		return "", false, err
	}
	raw, err = backend.Complete(system, prompt+feedback)
	if err == nil {
		return raw, false, nil
	}
	verdict := backendErrorVerdict(err)
	if aerr := applyVerdict(s, it, verdict, outPath, sessionPath); aerr != nil {
		return "", false, aerr
	}
	renderVerdict(out, it.Key, it, verdict)
	return "", true, nil
}
