//ff:func feature=cli type=helper control=sequence level=error
//ff:what scan이 def.Seed 에러를 표면화하는지 검증한다.

package cli

import (
	"path/filepath"
	"testing"
)

// TestScanSeedError: scan surfaces a def.Seed error.
func TestScanSeedError(t *testing.T) {
	dir := t.TempDir()
	session := filepath.Join(dir, "session.json")
	out := filepath.Join(dir, "out.jsonl")
	runCmdErr(t, errDef{seedErr: true}, session, out, "", "scan", "a")
}
