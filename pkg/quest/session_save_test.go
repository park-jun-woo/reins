//ff:func feature=quest type=helper control=sequence
//ff:what Save가 세션을 들여쓴 JSON으로 파일에 쓰는지 검증한다.

package quest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSaveWritesIndentedJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "session.json")
	s := &Session{Version: 1, Items: []*Item{{Key: "a", State: TODO}}}
	if err := s.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if !strings.Contains(string(b), "\n  ") {
		t.Errorf("output not indented: %q", b)
	}
}
