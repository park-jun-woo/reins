//ff:func feature=cli type=helper control=sequence
//ff:what NewQuestCmd가 빈 Options.Out에서 export 경로를 "<name>-results.jsonl"로 기본 설정하는지 검증한다.

package cli

import "testing"

// TestNewQuestCmdDefaultOut: with an empty Options.Out, the export path defaults to
// "<name>-results.jsonl".
func TestNewQuestCmdDefaultOut(t *testing.T) {
	cmd := NewQuestCmd("stub", stubDef{}, Options{})
	if got, _ := cmd.PersistentFlags().GetString("out"); got != "stub-results.jsonl" {
		t.Fatalf("default --out = %q, want stub-results.jsonl", got)
	}
}
