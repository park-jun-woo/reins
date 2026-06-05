//ff:func feature=cli type=helper control=sequence
//ff:what NewQuestCmd가 비지 않은 Options.Out으로 기본값을 덮어쓰고 Options.Version으로 루트 버전을 설정하는지 검증한다.

package cli

import "testing"

// TestNewQuestCmdExplicitOut: a non-empty Options.Out overrides the derived default,
// and Options.Version sets the root command's version.
func TestNewQuestCmdExplicitOut(t *testing.T) {
	cmd := NewQuestCmd("stub", stubDef{}, Options{Out: "custom.jsonl", Version: "9.9.9"})
	if got, _ := cmd.PersistentFlags().GetString("out"); got != "custom.jsonl" {
		t.Fatalf("default --out = %q, want custom.jsonl", got)
	}
	if cmd.Version != "9.9.9" {
		t.Fatalf("Version = %q, want 9.9.9", cmd.Version)
	}
}
