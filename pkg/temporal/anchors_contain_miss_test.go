//ff:func feature=temporal type=helper control=sequence
//ff:what anchorsContain이 토큰이 어떤 anchor에도 없으면 false를 내는지 검증한다.

package temporal

import "testing"

func TestAnchorsContainMiss(t *testing.T) {
	if anchorsContain([]string{"in Davos"}, "2099-12-31") {
		t.Fatal("token absent from all anchors should not be found")
	}
}
