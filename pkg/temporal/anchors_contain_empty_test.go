//ff:func feature=temporal type=helper control=sequence
//ff:what anchorsContain이 anchors가 비면 어떤 토큰도 포함하지 않아 false를 내는지 검증한다.

package temporal

import "testing"

func TestAnchorsContainEmpty(t *testing.T) {
	if anchorsContain(nil, "anything") {
		t.Fatal("no anchors should never contain a token")
	}
}
