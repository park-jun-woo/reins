//ff:func feature=cli type=helper control=sequence level=error
//ff:what scan이 세션 경로가 쓸 수 없을 때(부모 디렉터리 부재) save 에러를 표면화하는지 검증한다.

package cli

import (
	"path/filepath"
	"testing"
)

// TestScanSaveError: scan surfaces a save error when the session path is unwritable
// (parent directory absent).
func TestScanSaveError(t *testing.T) {
	dir := t.TempDir()
	session := filepath.Join(dir, "nope", "session.json")
	out := filepath.Join(dir, "out.jsonl")
	runCmdErr(t, stubDef{}, session, out, "", "scan", "a")
}
