//ff:func feature=cli type=command control=sequence level=error
//ff:what `scan [args...]` 명령. def.Seed로 입력에서 퀘스트 아이템을 시드하되 세션에 이미 있는 Key는 skip(dedupe — 중복 scan이 아이템을 복제해 래칫 무결성을 깨지 않게)하고 저장한다. 같은 입력 재scan은 신규분만 추가. skip이 있으면 "(skipped M duplicate(s))"를 병기.

package cli

import (
	"fmt"

	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/spf13/cobra"
)

// newScanCmd seeds new TODO items from args and saves the session. Seeded items
// whose Key already exists in the session are skipped (dedupe), so re-scanning the
// same input only adds the new ones and never duplicates the ratchet.
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
			existing := make(map[string]bool, len(s.Items))
			for _, it := range s.Items {
				existing[it.Key] = true
			}
			added, skipped := 0, 0
			for _, it := range items {
				if existing[it.Key] {
					skipped++
					continue
				}
				existing[it.Key] = true
				s.Items = append(s.Items, it)
				added++
			}
			if err := s.Save(*sessionPath); err != nil {
				return err
			}
			if skipped > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "seeded %d item(s) (skipped %d duplicate(s)); %d total\n", added, skipped, len(s.Items))
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "seeded %d item(s); %d total\n", added, len(s.Items))
			}
			return nil
		},
	}
}
