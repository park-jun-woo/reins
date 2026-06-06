//ff:func feature=cli type=command control=sequence level=error
//ff:what `agent [--model ...] [--max-items N]` 명령. submit의 자동 반복 — 남은 TODO를 NextTODO로 순회하며 LLM 생성(L0)→evaluateAndApply(게이트 판정·래칫·export, submit과 동일 경로)→FAIL이면 renderVerdictText 피드백을 user에, RuleSystem[verdict.RootCause] 코칭을 system에 되먹여 재시도(it.Tries<MaxTries). PASS/REVIEW/SKIP/BLOCK은 잠금→다음 아이템. backend는 opts.LLM!=nil이면 그걸, 아니면 --model을 llm.FromFlag로 lazy 생성. 종료 보장: MaxTries 초과 시 Apply가 DONE으로 잠가 NextTODO에서 빠짐(단조 수렴). PASS 잠금 권한은 게이트에만.

package cli

import (
	"fmt"

	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/llm"
	"github.com/park-jun-woo/reins/pkg/quest"
	"github.com/spf13/cobra"
)

// defaultAgentModel is the --model fallback when neither opts.DefaultModel nor the
// flag is set.
const defaultAgentModel = "ollama:gemma4:e4b"

// newAgentCmd builds the `agent` command: an automatic submit loop that lets the LLM
// generate each remaining TODO's payload, runs it through the same gate path as
// submit, and feeds FAIL feedback (plus rule-specific system coaching) back on retry.
func newAgentCmd(def gate.Definition, opts *AgentOptions, sessionPath, outPath *string, load sessionLoader) *cobra.Command {
	defaultModel := opts.DefaultModel
	if defaultModel == "" {
		defaultModel = defaultAgentModel
	}
	var (
		model    string
		maxItems int
	)
	cmd := &cobra.Command{
		Use:   "agent [--model backend:model] [--max-items N]",
		Short: "auto-run the generate→gate→retry loop over remaining TODO items",
		RunE: func(cmd *cobra.Command, args []string) error {
			backend := opts.LLM
			if backend == nil {
				b, err := llm.FromFlag(model)
				if err != nil {
					return err
				}
				backend = b
			}
			s, err := load()
			if err != nil {
				return err
			}
			// Signal Definition.Render to suppress its own last-failure log-tail
			// while the agent runs (the agent appends renderVerdict feedback itself,
			// avoiding double exposure). Cleared after the loop so a later manual
			// next/submit shows the tail again.
			s.SetMeta(quest.MetaAgentLoop, true)
			defer func() {
				delete(s.Meta, quest.MetaAgentLoop)
				_ = s.Save(*sessionPath)
			}()
			out := cmd.OutOrStdout()
			done := 0
			for it := s.NextTODO(); it != nil; it = s.NextTODO() {
				if maxItems > 0 && done >= maxItems {
					break
				}
				if err := runAgentItem(def, opts, backend, s, it, *outPath, *sessionPath, out); err != nil {
					return err
				}
				done++
			}
			fmt.Fprintf(out, "agent: processed %d item(s)\n", done)
			return nil
		},
	}
	cmd.Flags().StringVar(&model, "model", defaultModel, "LLM backend:model (ignored if a backend is injected)")
	cmd.Flags().IntVar(&maxItems, "max-items", 0, "max TODO items to process (0 = all)")
	return cmd
}
