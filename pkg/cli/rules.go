//ff:func feature=cli type=command control=sequence level=error
//ff:what `rules` 명령. gate.Catalog(def.Rules())로 게이트 규칙 카탈로그(레벨·ID·설명)를 출력한다 — 치즈 방어 감사를 위한 자동 rulebook.

package cli

import (
	"fmt"

	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/spf13/cobra"
)

// newRulesCmd prints the gate's rule catalog — the auto rulebook.
func newRulesCmd(def gate.Definition) *cobra.Command {
	return &cobra.Command{
		Use:   "rules",
		Short: "print the gate's rule catalog (auto rulebook)",
		RunE: func(cmd *cobra.Command, args []string) error {
			out := cmd.OutOrStdout()
			for _, m := range gate.Catalog(def.Rules()) {
				fmt.Fprintf(out, "%-6s %-24s %s\n", m.Level, m.ID, m.Desc)
			}
			return nil
		},
	}
}
