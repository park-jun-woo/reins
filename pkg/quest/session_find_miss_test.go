//ff:func feature=quest type=helper control=sequence level=error
//ff:what Find가 없는 Key에 대해 에러를 내는지 검증한다.

package quest

import "testing"

func TestFindMiss(t *testing.T) {
	s := &Session{Items: []*Item{{Key: "a"}}}
	if _, err := s.Find("zzz"); err == nil {
		t.Fatal("Find missing key: want error, got nil")
	}
}
