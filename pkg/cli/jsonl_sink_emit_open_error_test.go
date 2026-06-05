//ff:func feature=cli type=helper control=sequence level=error
//ff:what Emit이 파일을 열 수 없을 때(경로가 생성 후 디렉터리로 점유됨) 에러를 내는지 검증한다.

package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/reins/pkg/quest"
)

// TestEmitOpenError: Emit returns an error when its file cannot be opened (the path
// was usurped by a directory after construction).
func TestEmitOpenError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "occupied")
	if err := os.Mkdir(path, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	sink := &jsonlSink{path: path}
	if err := sink.Emit(&quest.Item{Key: "a", State: quest.PASS}); err == nil {
		t.Fatal("Emit to a directory path: want error, got nil")
	}
}
