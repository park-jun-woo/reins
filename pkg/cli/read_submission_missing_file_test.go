//ff:func feature=cli type=helper control=sequence level=error
//ff:what readSubmission이 부재 경로에서 read 에러를 표면화하는지 검증한다.

package cli

import (
	"path/filepath"
	"testing"
)

// TestReadSubmissionMissingFile: a nonexistent path surfaces the read error.
func TestReadSubmissionMissingFile(t *testing.T) {
	if _, err := readSubmission(newReadCmd(""), filepath.Join(t.TempDir(), "absent.txt")); err == nil {
		t.Fatal("readSubmission missing file: want error, got nil")
	}
}
