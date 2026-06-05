//ff:func feature=cli type=helper control=iteration dimension=1
//ff:what NewQuestCmd가 opts.ExtraCommands의 소비자 서브명령을 루트에 부착하는지 검증한다(G1).

package cli

import (
	"testing"

	"github.com/spf13/cobra"
)

// TestNewQuestCmdAttachesExtraCommands: a consumer command supplied via
// Options.ExtraCommands appears on the root alongside the canonical subcommands.
func TestNewQuestCmdAttachesExtraCommands(t *testing.T) {
	run := &cobra.Command{Use: "run", Short: "consumer run command"}
	cmd := NewQuestCmd("stub", stubDef{}, Options{ExtraCommands: []*cobra.Command{run}})

	var foundRun bool
	for _, c := range cmd.Commands() {
		if c.Name() == "run" {
			foundRun = true
		}
	}
	if !foundRun {
		t.Fatalf("ExtraCommands 'run' not attached to root")
	}
}
