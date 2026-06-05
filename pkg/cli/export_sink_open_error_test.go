//ff:func feature=cli type=helper control=sequence level=error
//ff:what export가 out 경로의 부모를 만들 수 없을 때(부모 슬롯을 일반 파일이 차지) sink 생성 에러를 표면화하는지 검증한다.

package cli

import (
	"os"
	"path/filepath"
	"testing"
)

// TestExportSinkOpenError: export surfaces a sink-construction error when the out
// path's parent cannot be created (a regular file occupies a parent slot).
func TestExportSinkOpenError(t *testing.T) {
	dir := t.TempDir()
	session := filepath.Join(dir, "session.json")
	blocker := filepath.Join(dir, "blocker")
	if err := os.WriteFile(blocker, []byte("x"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	out := filepath.Join(blocker, "sub", "out.jsonl")
	runCmdErr(t, stubDef{}, session, out, "", "export")
}
