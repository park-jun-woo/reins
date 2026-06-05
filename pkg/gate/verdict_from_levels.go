//ff:func feature=gate type=helper control=selection
//ff:what 집계된 레벨 플래그를 Verdict로 환원한다. anyFail이면 FAIL(+Facts), 아니면 anyReview면 REVIEW(+Facts), 아니면 PASS. 가중치 아님 — 결정적 위반 1개가 곧 FAIL.

package gate

import "github.com/park-jun-woo/reins/pkg/quest"

// verdictFromLevels reduces the aggregated level flags to a Verdict: any Fail → FAIL
// (with Facts), else any Review → REVIEW (with Facts), else PASS. Never a weight — a
// single decisive violation is FAIL.
func verdictFromLevels(anyFail, anyReview bool, facts []quest.Fact) quest.Verdict {
	switch {
	case anyFail:
		return quest.Verdict{Outcome: quest.OutFail, Facts: facts}
	case anyReview:
		return quest.Verdict{Outcome: quest.OutReview, Facts: facts}
	default:
		return quest.Verdict{Outcome: quest.OutPass}
	}
}
