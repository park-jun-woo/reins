//ff:func feature=cli type=helper control=sequence level=error
//ff:what status가 깨진 세션 파일의 load 에러를 표면화하는지 검증한다.

package cli

import (
	"os"
	"path/filepath"
	"testing"
)

// TestStatusLoadError: status surfaces a load error from a corrupt session file.
func TestStatusLoadError(t *testing.T) {
	dir := t.TempDir()
	session := filepath.Join(dir, "session.json")
	out := filepath.Join(dir, "out.jsonl")
	if err := os.WriteFile(session, []byte("{bad"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	runCmdErr(t, stubDef{}, session, out, "", "status")
}
