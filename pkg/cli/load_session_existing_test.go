//ff:func feature=cli type=helper control=sequence level=error
//ff:what loadSession이 기존 유효 세션 파일을 그대로 로드하는지 검증한다.

package cli

import (
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/reins/pkg/quest"
)

// TestLoadSessionExisting: an existing valid session file is loaded as-is.
func TestLoadSessionExisting(t *testing.T) {
	path := filepath.Join(t.TempDir(), "session.json")
	want := &quest.Session{Version: 1, Items: []*quest.Item{{Key: "a", State: quest.TODO}}}
	if err := want.Save(path); err != nil {
		t.Fatalf("save: %v", err)
	}
	s, err := loadSession(path)
	if err != nil {
		t.Fatalf("loadSession: %v", err)
	}
	if len(s.Items) != 1 || s.Items[0].Key != "a" {
		t.Fatalf("loadSession = %+v, want loaded items", s)
	}
}
