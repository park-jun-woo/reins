//ff:func feature=graph type=helper control=iteration dimension=1
//ff:what verdictFromCounters 증명 — 레벨 플래그를 Verdict로 환원: anyFail→FAIL(+Facts), anyFail 없고 anyReview→REVIEW(+Facts), 둘 다 없으면 PASS(Facts 없음). Fail 우선순위(둘 다 true면 FAIL)를 테이블로 못 박는다.

package graph

import (
	"testing"

	"github.com/park-jun-woo/reins/pkg/quest"
)

func TestVerdictFromCounters(t *testing.T) {
	facts := []quest.Fact{{Where: "x"}}

	cases := []struct {
		name        string
		anyFail     bool
		anyReview   bool
		wantOutcome quest.Outcome
		wantFacts   bool
	}{
		{"fail", true, false, quest.OutFail, true},
		{"review", false, true, quest.OutReview, true},
		{"fail dominates review", true, true, quest.OutFail, true},
		{"pass", false, false, quest.OutPass, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			runVerdictFromCountersCase(t, c.anyFail, c.anyReview, facts, c.wantOutcome, c.wantFacts)
		})
	}
}
