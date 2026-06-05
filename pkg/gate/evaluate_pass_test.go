//ff:func feature=gate type=helper control=sequence
//ff:what 아무 규칙도 발동하지 않으면 Evaluate가 OutPass를 내는지 검증한다.

package gate

import (
	"testing"

	"github.com/park-jun-woo/reins/pkg/quest"
)

func TestEvaluate_NoRuleFires_Pass(t *testing.T) {
	rules := []Rule{fireRule("a", LevelFail, false), fireRule("b", LevelReview, false)}
	v := Evaluate(rules, Context{})
	if v.Outcome != quest.OutPass {
		t.Fatalf("outcome = %s, want %s", v.Outcome, quest.OutPass)
	}
}
