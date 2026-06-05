//ff:func feature=gate type=helper control=sequence
//ff:what FAIL 규칙이 발동하면 REVIEW 규칙이 함께 발동해도 Evaluate가 OutFail을 내는지(레벨 집계에서 FAIL 우선) 검증한다.

package gate

import (
	"testing"

	"github.com/park-jun-woo/reins/pkg/quest"
)

func TestEvaluate_FailFires_FailWinsOverReview(t *testing.T) {
	rules := []Rule{fireRule("a", LevelFail, true), fireRule("b", LevelReview, true)}
	v := Evaluate(rules, Context{})
	if v.Outcome != quest.OutFail {
		t.Fatalf("outcome = %s, want %s", v.Outcome, quest.OutFail)
	}
}
