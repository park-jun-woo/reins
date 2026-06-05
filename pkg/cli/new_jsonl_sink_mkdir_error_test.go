//ff:func feature=cli type=helper control=sequence level=error
//ff:what newJSONLSink이 부모 슬롯을 일반 파일이 차지하면 MkdirAll 에러를 표면화하는지 검증한다.

package cli

import (
	"os"
	"path/filepath"
	"testing"
)

// TestNewJSONLSinkMkdirError: newJSONLSink surfaces a MkdirAll error when a regular
// file occupies a parent slot of the requested nested path.
func TestNewJSONLSinkMkdirError(t *testing.T) {
	dir := t.TempDir()
	blocker := filepath.Join(dir, "blocker")
	if err := os.WriteFile(blocker, []byte("x"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	if _, err := newJSONLSink(filepath.Join(blocker, "sub", "out.jsonl")); err == nil {
		t.Fatal("newJSONLSink under a file: want error, got nil")
	}
}
