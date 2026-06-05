//ff:func feature=quest type=helper control=sequence
//ff:what New가 스키마 버전 1의 빈(아이템 0개) 세션을 새로 만드는지 검증한다.

package quest

import "testing"

func TestNewEmptySession(t *testing.T) {
	s := New()
	if s == nil {
		t.Fatal("New returned nil")
	}
	if s.Version != 1 {
		t.Errorf("Version = %d, want 1", s.Version)
	}
	if len(s.Items) != 0 {
		t.Errorf("Items = %d, want empty", len(s.Items))
	}
}
