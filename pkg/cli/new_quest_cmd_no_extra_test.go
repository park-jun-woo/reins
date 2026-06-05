//ff:func feature=cli type=helper control=sequence
//ff:what nil ExtraCommands가 표준 6개 서브명령만 남기는지 검증한다(G1 후방호환).

package cli

import "testing"

// TestNewQuestCmdNoExtraCommands: nil ExtraCommands leaves only the canonical
// subcommands (backward-compatible).
func TestNewQuestCmdNoExtraCommands(t *testing.T) {
	cmd := NewQuestCmd("stub", stubDef{}, Options{})
	if got := len(cmd.Commands()); got != 6 {
		t.Fatalf("subcommand count = %d, want 6 canonical with no extras", got)
	}
}
