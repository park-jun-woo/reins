//ff:func feature=temporal type=helper control=sequence
//ff:what 성분(라틴숫자 런)이 anchor에 실재하면 ComponentsInAnchor가 true를 내는지 검증한다.

package temporal

import "testing"

func TestComponentsInAnchorPresent(t *testing.T) {
	spec := Spec{Kind: Absolute, Calendar: Gregorian, Start: "2017-01-10"}
	if !ComponentsInAnchor(spec, []string{"on 2017-01-10 in Davos"}) {
		t.Fatal("components present should pass")
	}
}
