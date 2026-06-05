//ff:func feature=cli type=helper control=sequence level=error
//ff:what loadSession이 깨진(부재 아님) 세션 파일에서 load 에러를 표면화하는지(빈 세션으로 가리지 않는지) 검증한다.

package cli

import (
	"os"
	"path/filepath"
	"testing"
)

// TestLoadSessionBadJSON: a corrupt (non-missing) session file surfaces the load
// error rather than masking it as a fresh session.
func TestLoadSessionBadJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "broken.json")
	if err := os.WriteFile(path, []byte("{bad"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	if _, err := loadSession(path); err == nil {
		t.Fatal("loadSession broken JSON: want error, got nil")
	}
}
