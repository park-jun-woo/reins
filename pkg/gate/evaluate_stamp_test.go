//ff:func feature=gate type=helper control=sequence
//ff:what 발동한 규칙의 Fact가 그 규칙 ID로 스탬프되어 verdict의 evidence로 수집되는지 검증한다.

package gate

import (
	"testing"
)

func TestEvaluate_StampsFiredRuleID(t *testing.T) {
	v := Evaluate([]Rule{fireRule("r1", LevelFail, true)}, Context{})
	if len(v.Facts) != 1 || v.Facts[0].Rule != "r1" {
		t.Fatalf("facts = %+v, want one fact stamped with r1", v.Facts)
	}
}
