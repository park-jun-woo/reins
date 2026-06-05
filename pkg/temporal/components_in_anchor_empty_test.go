//ff:func feature=temporal type=helper control=sequence
//ff:what 라틴숫자 성분이 0개면(예: Relative 명세) 묶을 게 없어 ComponentsInAnchor가 true를 내는지 검증한다.

package temporal

import "testing"

func TestComponentsInAnchorEmpty(t *testing.T) {
	if !ComponentsInAnchor(Spec{Kind: Relative}, nil) {
		t.Fatal("no components ⇒ nothing to tie ⇒ true")
	}
}
