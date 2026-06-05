//ff:func feature=graph type=helper control=selection
//ff:what verdictFromCounters — 잔존 카운터의 레벨 플래그를 Verdict로 환원한다. anyFail→FAIL(+Facts), 아니면 anyReview→REVIEW(+Facts), 아니면 PASS(Facts 없음). gate.verdictFromLevels와 동일 의미 — 가중치 아님, 잔존 Fail 1개가 곧 FAIL.

package graph

import "github.com/park-jun-woo/reins/pkg/quest"

// verdictFromCounters reduces the remaining counters' level flags to a Verdict:
// any Fail ⇒ FAIL (with Facts), else any Review ⇒ REVIEW (with Facts), else PASS
// (no Facts). Same semantics as gate.verdictFromLevels — never a weight; a single
// remaining Fail counter is FAIL.
func verdictFromCounters(anyFail, anyReview bool, facts []quest.Fact) quest.Verdict {
	switch {
	case anyFail:
		return quest.Verdict{Outcome: quest.OutFail, Facts: facts}
	case anyReview:
		return quest.Verdict{Outcome: quest.OutReview, Facts: facts}
	default:
		return quest.Verdict{Outcome: quest.OutPass}
	}
}
