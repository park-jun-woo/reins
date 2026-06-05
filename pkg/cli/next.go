//ff:func feature=cli type=command control=sequence level=error
//ff:what `next` 명령(읽기·비변이). 다음 TODO 아이템 하나를 골라 def.Render로 작성 프롬프트·검증 컨텍스트를 출력한다. TODO가 없으면 그 사실만 알린다.

package cli

import (
	"fmt"

	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/spf13/cobra"
)

// newNextCmd prints the next TODO item's authoring prompt + verification context.
func newNextCmd(def gate.Definition, load sessionLoader) *cobra.Command {
	return &cobra.Command{
		Use:   "next",
		Short: "show the next TODO item (read-only)",
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := load()
			if err != nil {
				return err
			}
			it := s.NextTODO()
			if it == nil {
				fmt.Fprintln(cmd.OutOrStdout(), "no TODO items remaining")
				return nil
			}
			prompt, err := def.Render(it)
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), prompt)
			return nil
		},
	}
}
