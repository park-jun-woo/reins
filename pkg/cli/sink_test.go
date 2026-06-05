//ff:func feature=cli type=helper control=sequence level=error
//ff:what newJSONLSink이 중첩 out 경로의 부모 디렉터리를 만드는지(첫 Emit이 파일을 열 수 있게) 검증한다.

package cli

import (
	"os"
	"path/filepath"
	"testing"
)

// TestNewJSONLSinkCreatesDir: newJSONLSink creates the parent directory for a nested
// out path so the first Emit can open the file.
func TestNewJSONLSinkCreatesDir(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "out.jsonl")
	sink, err := newJSONLSink(path)
	if err != nil {
		t.Fatalf("newJSONLSink: %v", err)
	}
	if _, err := os.Stat(filepath.Dir(path)); err != nil {
		t.Fatalf("parent dir not created: %v", err)
	}
	if sink.path != path {
		t.Fatalf("sink.path = %q, want %q", sink.path, path)
	}
}
