//ff:func feature=cli type=helper control=sequence level=error
//ff:what readSubmission이 구체 경로에서 파일 바이트를 읽는지 검증한다.

package cli

import (
	"os"
	"path/filepath"
	"testing"
)

// TestReadSubmissionFile: a concrete path reads the file's bytes.
func TestReadSubmissionFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sub.txt")
	if err := os.WriteFile(path, []byte("from file"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	got, err := readSubmission(newReadCmd(""), path)
	if err != nil {
		t.Fatalf("readSubmission: %v", err)
	}
	if string(got) != "from file" {
		t.Fatalf("readSubmission = %q, want file content", got)
	}
}
