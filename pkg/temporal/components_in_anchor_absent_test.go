//ff:func feature=temporal type=helper control=sequence
//ff:what 성분이 anchor에 없으면(엉뚱앵커 Türkiye·Kurtulmuş류) ComponentsInAnchor가 false로 차단하는지 검증한다.

package temporal

import "testing"

func TestComponentsInAnchorAbsent(t *testing.T) {
	spec := Spec{Kind: Absolute, Calendar: Gregorian, Start: "2017-01-10"}
	if ComponentsInAnchor(spec, []string{"Türkiye Kurtulmuş said"}) {
		t.Fatal("components absent should fail")
	}
}
