//ff:func feature=cli type=helper control=iteration dimension=1
//ff:what NewQuestCmd가 표준 서브커맨드(scan/next/submit/status/export/rules)를 모두 배선하는지 검증한다.

package cli

import "testing"

// TestNewQuestCmdWiresSubcommands: the root wires the canonical subcommands.
func TestNewQuestCmdWiresSubcommands(t *testing.T) {
	cmd := NewQuestCmd("stub", stubDef{}, Options{})
	want := map[string]bool{"scan": false, "next": false, "submit": false, "status": false, "export": false, "rules": false}
	for _, c := range cmd.Commands() {
		want[c.Name()] = true
	}
	for name, found := range want {
		if !found {
			t.Errorf("subcommand %q not wired", name)
		}
	}
}
