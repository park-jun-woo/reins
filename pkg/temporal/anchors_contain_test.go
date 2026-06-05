//ff:func feature=temporal type=helper control=sequence
//ff:what anchorsContain이 토큰이 anchors 중 하나에라도 substring으로 들어있으면 true를 내는지 검증한다.

package temporal

import "testing"

func TestAnchorsContainHit(t *testing.T) {
	anchors := []string{"in Davos", "on 2017-01-10"}
	if !anchorsContain(anchors, "2017-01-10") {
		t.Fatal("token present in an anchor should be found")
	}
}
