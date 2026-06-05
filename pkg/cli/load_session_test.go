//ff:func feature=cli type=helper control=sequence level=error
//ff:what loadSession이 부재 세션 파일에서 에러 대신 빈 새 세션을 내는지 검증한다(첫 scan을 위해).

package cli

import (
	"path/filepath"
	"testing"
)

// TestLoadSessionAbsentFresh: a missing session file yields a fresh empty session
// rather than an error, so the first scan can start from nothing.
func TestLoadSessionAbsentFresh(t *testing.T) {
	s, err := loadSession(filepath.Join(t.TempDir(), "absent.json"))
	if err != nil {
		t.Fatalf("loadSession absent: %v", err)
	}
	if s == nil || s.Version != 1 || len(s.Items) != 0 {
		t.Fatalf("loadSession absent = %+v, want fresh quest.New()", s)
	}
}
