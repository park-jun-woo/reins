//ff:func feature=gate type=helper control=sequence
//ff:what Catalog가 입력 규칙들의 RuleMeta를 순서대로 반환하는지 검증한다(자동 rulebook).

package gate

import (
	"testing"
)

func TestCatalog_ReturnsRuleMetas(t *testing.T) {
	metas := Catalog([]Rule{fireRule("a", LevelFail, false), fireRule("b", LevelReview, false)})
	if len(metas) != 2 || metas[0].ID != "a" || metas[1].ID != "b" {
		t.Fatalf("catalog = %+v", metas)
	}
}
