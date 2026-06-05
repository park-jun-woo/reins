//ff:func feature=cli type=command control=sequence level=error
//ff:what `status` 명령. Session.Progress로 상태별 집계(TODO/PASS/REVIEW/DONE/SKIPPED/BLOCKED)와 TOTAL을 출력한다.

package cli

import (
	"fmt"

	"github.com/park-jun-woo/reins/pkg/quest"
	"github.com/spf13/cobra"
)

// newStatusCmd prints the per-state tally from the session.
func newStatusCmd(load sessionLoader) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "show the progress tally",
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := load()
			if err != nil {
				return err
			}
			out := cmd.OutOrStdout()
			prog := s.Progress()
			for _, st := range []quest.State{quest.TODO, quest.PASS, quest.REVIEW, quest.DONE, quest.SKIPPED, quest.BLOCKED} {
				fmt.Fprintf(out, "%-8s %d\n", st, prog[st])
			}
			fmt.Fprintf(out, "%-8s %d\n", "TOTAL", len(s.Items))
			return nil
		},
	}
}
