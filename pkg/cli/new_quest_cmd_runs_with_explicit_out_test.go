//ff:func feature=cli type=helper control=sequence level=error
//ff:what NewQuestCmd의 명시 out 경로가 end-to-end로 실제 사용되는지(export가 거기 쓰는지) 검증한다.

package cli

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

// TestNewQuestCmdRunsWithExplicitOut: the explicit out path is actually used end to
// end (export writes to it).
func TestNewQuestCmdRunsWithExplicitOut(t *testing.T) {
	dir := t.TempDir()
	session := filepath.Join(dir, "session.json")
	out := filepath.Join(dir, "explicit.jsonl")

	cmd := NewQuestCmd("stub", stubDef{}, Options{Out: out})
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetIn(strings.NewReader(""))
	cmd.SetArgs([]string{"--session", session, "scan", "a"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("scan: %v\n%s", err, buf.String())
	}
}
