//ff:func feature=cli type=helper control=iteration dimension=1
//ff:what runAgentItem — 한 TODO 아이템의 생성→게이트→재시도 루프. it.State==TODO이고 Tries<MaxTries인 동안 composeSystem(피드백 코칭)으로 system을 짜고 def.Render→backend.Complete로 페이로드를 생성, evaluateAndApply로 게이트 판정·래칫·export 후 renderVerdict 출력. FAIL이 아니면 종료, FAIL이면 renderVerdictText 피드백을 user에, RuleSystem[verdict.RootCause] 코칭을 다음 system에 되먹인다. MaxTries 초과 시 Apply가 DONE으로 잠가 루프 탈출(단조 수렴).

package cli

import (
	"io"

	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/llm"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// runAgentItem runs the generate→gate→retry loop for one TODO item, feeding FAIL
// feedback (and rule-specific system coaching) back on each retry until the item
// leaves TODO or exhausts MaxTries.
func runAgentItem(def gate.Definition, opts *AgentOptions, backend llm.Backend, s *quest.Session, it *quest.Item, outPath, sessionPath string, out io.Writer) error {
	feedback, ruleCoach := "", ""
	for it.State == quest.TODO && it.Tries < quest.MaxTries {
		system := composeSystem(opts.System, ruleCoach)
		prompt, err := def.Render(s, it)
		if err != nil {
			return err
		}
		user := prompt + feedback
		raw, err := backend.Complete(system, user)
		if err != nil {
			return err
		}
		verdict, err := evaluateAndApply(def, s, it, []byte(raw), outPath, sessionPath)
		if err != nil {
			return err
		}
		renderVerdict(out, it.Key, it, verdict)
		if verdict.Outcome != quest.OutFail {
			break
		}
		feedback = "\n\n--- PREVIOUS ATTEMPT FAILED ---\n" + renderVerdictText(it.Key, it, verdict)
		ruleCoach = opts.RuleSystem[verdict.RootCause]
	}
	return nil
}
