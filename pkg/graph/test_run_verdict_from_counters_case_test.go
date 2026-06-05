//ff:func feature=graph type=helper control=sequence
//ff:what runVerdictFromCountersCase 테스트헬퍼 — 레벨 플래그로 verdictFromCounters를 호출해 Outcome과 Facts 동반 여부(PASS는 무Facts)를 단언한다(케이스 1건 판정).

package graph

import (
	"testing"

	"github.com/park-jun-woo/reins/pkg/quest"
)

func runVerdictFromCountersCase(t *testing.T, anyFail, anyReview bool, facts []quest.Fact, wantOutcome quest.Outcome, wantFacts bool) {
	t.Helper()
	v := verdictFromCounters(anyFail, anyReview, facts)
	if v.Outcome != wantOutcome {
		t.Fatalf("outcome=%s want=%s", v.Outcome, wantOutcome)
	}
	if wantFacts && len(v.Facts) != 1 {
		t.Fatalf("expected facts carried, got %v", v.Facts)
	}
	if !wantFacts && len(v.Facts) != 0 {
		t.Fatalf("PASS should carry no facts, got %v", v.Facts)
	}
}
