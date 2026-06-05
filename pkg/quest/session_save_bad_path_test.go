//ff:func feature=quest type=helper control=sequence level=error
//ff:what Save가 쓸 수 없는 경로(부모 디렉터리 부재)에선 에러를 내는지 검증한다.

package quest

import (
	"path/filepath"
	"testing"
)

func TestSaveBadPath(t *testing.T) {
	s := &Session{Version: 1}
	// Parent directory does not exist, so WriteFile fails.
	path := filepath.Join(t.TempDir(), "nope", "session.json")
	if err := s.Save(path); err == nil {
		t.Fatal("Save to nonexistent dir: want error, got nil")
	}
}
