//ff:func feature=quest type=helper control=sequence
//ff:what NextTODO가 첫 TODO 아이템을 고르고 잠긴 것은 건너뛰는지 검증한다(래칫).

package quest

import "testing"

func TestNextTODOSkipsLocked(t *testing.T) {
	s := &Session{Items: []*Item{
		{Key: "a", State: PASS}, // locked, skipped
		{Key: "b", State: TODO}, // first TODO
		{Key: "c", State: TODO}, // not reached
	}}
	it := s.NextTODO()
	if it == nil || it.Key != "b" {
		t.Fatalf("NextTODO = %+v, want item b", it)
	}
}
