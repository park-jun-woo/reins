//ff:func feature=quest type=helper control=sequence
//ff:what NextTODO가 TODO가 없으면 nil을 내는지 검증한다(잠긴 것뿐일 때).

package quest

import "testing"

func TestNextTODONone(t *testing.T) {
	s := &Session{Items: []*Item{{Key: "a", State: DONE}}}
	if it := s.NextTODO(); it != nil {
		t.Fatalf("NextTODO = %+v, want nil when nothing remains", it)
	}
}
