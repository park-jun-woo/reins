//ff:func feature=cli type=helper control=iteration dimension=1
//ff:what runLoopItem — 한 TODO 아이템의 생성→게이트→재시도 루프. it.State==TODO이고 Tries<MaxTries인 동안 composeSystem(피드백 코칭)으로 system을 짜고 def.Render→backend.Complete로 페이로드를 생성, evaluateAndApply로 게이트 판정·래칫·export 후 renderVerdict 출력. backend.Complete 에러는 handleBackendError로 강등해 backendErrorVerdict를 래칫(Tries++)·계속(루프 abort 아님; 영속화 실패만 전파). FAIL이 아니면 종료, FAIL이면 renderVerdictText 피드백을 user에, RuleSystem[verdict.RootCause] 코칭을 다음 system에 되먹인다(생성 오류 FAIL은 내용 비평이 아니므로 피드백 비주입). FAIL의 RootCause가 opts.EscalateOn에 있고 opts.Escalate가 있으면 그 아이템을 강한 backend로 래치 승격(능력 한계 신호에만; 형식 실패는 약한 모델 유지). MaxTries 초과 시 Apply가 DONE으로 잠가 루프 탈출(단조 수렴) — 승격 후에도 실패면 최종 FAIL. def.Render·applyVerdict(Save/export) 에러는 그대로 전파(fail-fast).

package cli

import (
	"fmt"
	"io"

	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/llm"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// runLoopItem runs the generate→gate→retry loop for one TODO item, feeding FAIL
// feedback (and rule-specific system coaching) back on each retry until the item
// leaves TODO or exhausts MaxTries.
//
// Escalation (opts.Escalate + opts.EscalateOn): the loop starts on the primary
// backend. When a FAIL's RootCause is in EscalateOn — a capability-bound signal the
// consumer designates, not a format slip — the item latches onto the stronger
// Escalate backend for its remaining tries. The gate is still the sole PASS
// authority; escalation only swaps the L0 generator, so a stronger model that still
// FAILs simply exhausts MaxTries and locks DONE (the true residual).
func runLoopItem(def gate.Definition, opts *LoopOptions, primary llm.Backend, s *quest.Session, it *quest.Item, outPath, sessionPath string, out io.Writer) error {
	escalateOn := escalateRootCauses(opts.EscalateOn)
	escalated := false
	feedback, ruleCoach := "", ""
	for it.State == quest.TODO && it.Tries < quest.MaxTries {
		backend := primary
		if escalated && opts.Escalate != nil {
			backend = opts.Escalate
		}
		raw, handled, err := generatePayload(def, opts, backend, ruleCoach, feedback, s, it, outPath, sessionPath, out)
		if err != nil {
			return err
		}
		if handled {
			continue
		}
		verdict, err := evaluateAndApply(def, s, it, []byte(raw), outPath, sessionPath)
		if err != nil {
			return err
		}
		renderVerdict(out, it.Key, it, verdict)
		if verdict.Outcome != quest.OutFail {
			break
		}
		// Promote to the stronger backend when the primary proves capability-bound:
		// a FAIL whose RootCause the consumer marked escalation-worthy. Latch on —
		// once escalated, the item stays on Escalate for its remaining tries.
		if !escalated && opts.Escalate != nil && escalateOn[verdict.RootCause] {
			escalated = true
			fmt.Fprintf(out, "  ↑ escalating %s to stronger model (RootCause=%s)\n", it.Key, verdict.RootCause)
		}
		feedback = "\n\n--- PREVIOUS ATTEMPT FAILED ---\n" + renderVerdictText(it.Key, it, verdict)
		ruleCoach = opts.RuleSystem[verdict.RootCause]
	}
	return nil
}
