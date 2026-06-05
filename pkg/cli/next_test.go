//ff:func feature=cli type=helper control=sequence
//ff:what next가 첫 TODO 아이템을 def.Render로 렌더하는지 검증한다.

package cli

import (
	"path/filepath"
	"strings"
	"testing"
)

// TestNextRenders: next renders the first TODO item via def.Render.
func TestNextRenders(t *testing.T) {
	dir := t.TempDir()
	session := filepath.Join(dir, "session.json")
	out := filepath.Join(dir, "out.jsonl")
	runCmd(t, stubDef{}, session, out, "", "scan", "a", "b")

	got := runCmd(t, stubDef{}, session, out, "", "next")
	if !strings.Contains(got, "render:a") {
		t.Fatalf("next = %q", got)
	}
}
