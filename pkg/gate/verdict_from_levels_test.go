//ff:func feature=gate type=helper control=iteration dimension=1
//ff:what verdictFromLevels가 레벨 플래그를 Verdict로 환원하는지 검증한다 — anyFail→FAIL(우선), anyReview→REVIEW, 둘 다 거짓→PASS, Facts 전달.

package gate

import (
	"testing"

	"github.com/park-jun-woo/reins/pkg/quest"
)

func TestVerdictFromLevels(t *testing.T) {
	facts := []quest.Fact{{Rule: "r1", Where: "who"}}
	cases := []struct {
		name      string
		anyFail   bool
		anyReview bool
		want      quest.Outcome
		wantFacts bool
	}{
		{"fail wins over review", true, true, quest.OutFail, true},
		{"only fail", true, false, quest.OutFail, true},
		{"only review", false, true, quest.OutReview, true},
		{"none passes, no facts", false, false, quest.OutPass, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			v := verdictFromLevels(c.anyFail, c.anyReview, facts)
			if v.Outcome != c.want {
				t.Fatalf("outcome = %q, want %q", v.Outcome, c.want)
			}
			if got := len(v.Facts) > 0; got != c.wantFacts {
				t.Fatalf("hasFacts = %v, want %v", got, c.wantFacts)
			}
		})
	}
}
