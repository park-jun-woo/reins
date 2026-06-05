//ff:func feature=quest type=helper control=sequence
//ff:what Find가 Key로 아이템을 정확히 찾는지 검증한다.

package quest

import "testing"

func TestFindHit(t *testing.T) {
	want := &Item{Key: "b", State: TODO}
	s := &Session{Items: []*Item{{Key: "a"}, want}}
	got, err := s.Find("b")
	if err != nil {
		t.Fatalf("Find error: %v", err)
	}
	if got != want {
		t.Fatalf("Find = %+v, want %+v", got, want)
	}
}
