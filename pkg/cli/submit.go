//ff:func feature=cli type=command control=sequence level=error
//ff:what `submit --key <k> [--in <file>|-]` 명령. 세션 Find→TODO 확인→제출물 raw 로드→def.Prepare(s, it, raw)(세션을 넘겨 Prepare가 s.Meta를 읽고 갱신; 직후 Save로 영속)→(short verdict, def가 gate.Evaluator면 ev.Evaluate(그래프), 아니면 gate.Evaluate(Rules))→quest.Apply(UTC RFC3339)→Save→Export 후, 결과와 FAIL 시 Fact 피드백을 출력한다.

package cli

import (
	"fmt"
	"time"

	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
	"github.com/spf13/cobra"
)

// newSubmitCmd evaluates a submission for one item through the gate, applies the
// ratchet transition, exports terminal items, and prints the outcome plus any Facts.
func newSubmitCmd(def gate.Definition, sessionPath, outPath *string, load sessionLoader) *cobra.Command {
	var (
		key    string
		inPath string
	)
	cmd := &cobra.Command{
		Use:   "submit --key <k> [--in <file>|-]",
		Short: "submit an item for gate evaluation",
		RunE: func(cmd *cobra.Command, args []string) error {
			if key == "" {
				return fmt.Errorf("--key is required")
			}
			s, err := load()
			if err != nil {
				return err
			}
			it, err := s.Find(key)
			if err != nil {
				return err
			}
			if it.State != quest.TODO {
				return fmt.Errorf("item %s is %s, not TODO", key, it.State)
			}
			raw, err := readSubmission(cmd, inPath)
			if err != nil {
				return err
			}
			ctx, short, err := def.Prepare(s, it, raw)
			if err != nil {
				return err
			}
			var verdict quest.Verdict
			if short != nil {
				verdict = *short
			} else if ev, ok := def.(gate.Evaluator); ok {
				verdict = ev.Evaluate(ctx)
			} else {
				verdict = gate.Evaluate(def.Rules(), ctx)
			}
			now := time.Now().UTC().Format(time.RFC3339)
			quest.Apply(it, verdict, now)
			if err := s.Save(*sessionPath); err != nil {
				return err
			}
			sink, err := newJSONLSink(*outPath)
			if err != nil {
				return err
			}
			if _, err := quest.Export(s, sink); err != nil {
				return err
			}
			if err := s.Save(*sessionPath); err != nil {
				return err
			}
			printSubmit(cmd.OutOrStdout(), key, it, verdict)
			return nil
		},
	}
	cmd.Flags().StringVar(&key, "key", "", "item key to submit (required)")
	cmd.Flags().StringVar(&inPath, "in", "-", "submission file, or - for stdin")
	return cmd
}
