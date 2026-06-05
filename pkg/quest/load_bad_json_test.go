//ff:func feature=quest type=helper control=sequence level=error
//ff:what Load가 깨진 JSON에서 에러를 내는지 검증한다.

package quest

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadBadJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "broken.json")
	if err := os.WriteFile(path, []byte("{not json"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	if _, err := Load(path); err == nil {
		t.Fatal("Load broken JSON: want error, got nil")
	}
}
