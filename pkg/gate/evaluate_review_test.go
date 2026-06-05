//ff:func feature=gate type=helper control=sequence
//ff:what FAIL 규칙은 잠잠하고 REVIEW 규칙만 발동하면 Evaluate가 OutReview를 내는지 검증한다.

package gate

import (
	"testing"

	"github.com/park-jun-woo/reins/pkg/quest"
)

func TestEvaluate_OnlyReviewFires_Review(t *testing.T) {
	rules := []Rule{fireRule("a", LevelFail, false), fireRule("b", LevelReview, true)}
	v := Evaluate(rules, Context{})
	if v.Outcome != quest.OutReview {
		t.Fatalf("outcome = %s, want %s", v.Outcome, quest.OutReview)
	}
}
