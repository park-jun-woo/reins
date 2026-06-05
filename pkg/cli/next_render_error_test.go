//ff:func feature=cli type=helper control=sequence level=error
//ff:what next가 다음 TODO 아이템의 def.Render 에러를 표면화하는지 검증한다.

package cli

import (
	"path/filepath"
	"testing"
)

// TestNextRenderError: next surfaces a def.Render error for the next TODO item.
func TestNextRenderError(t *testing.T) {
	dir := t.TempDir()
	session := filepath.Join(dir, "session.json")
	out := filepath.Join(dir, "out.jsonl")
	runCmd(t, errDef{}, session, out, "", "scan", "a")
	runCmdErr(t, errDef{renderErr: true}, session, out, "", "next")
}
