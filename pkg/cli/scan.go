//ff:func feature=cli type=command control=sequence level=error
//ff:what `scan [args...]` 명령. def.Seed로 입력에서 N개 퀘스트 아이템을 시드해 세션에 추가하고 저장한다(scan은 입력에서 작업을 시드하는 본질).

package cli

import (
	"fmt"

	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/spf13/cobra"
)

// newScanCmd seeds new TODO items from args and saves the session.
func newScanCmd(def gate.Definition, sessionPath *string, load sessionLoader) *cobra.Command {
	return &cobra.Command{
		Use:   "scan [args...]",
		Short: "seed quest items from the input",
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := load()
			if err != nil {
				return err
			}
			items, err := def.Seed(args)
			if err != nil {
				return err
			}
			s.Items = append(s.Items, items...)
			if err := s.Save(*sessionPath); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "seeded %d item(s); %d total\n", len(items), len(s.Items))
			return nil
		},
	}
}
