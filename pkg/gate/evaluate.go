//ff:func feature=gate type=helper control=iteration dimension=1
//ff:what 카탈로그를 제출 1건에 평가→발동 규칙을 레벨로 집계(가중치 아님). 발동 Fact는 규칙 ID로 스탬프되어 evidence로 수집되고, 집계된 레벨은 verdictFromLevels가 Verdict로 환원한다. 결정론: 같은 (rules, ctx)→같은 Verdict. defeat/h-Categoriser는 toulmin 백엔드 도입 시 여기 플러그인.

package gate

import "github.com/park-jun-woo/reins/pkg/quest"

// Evaluate runs every rule's Check over one submission and aggregates the fired
// rules by Level (never by weight): any fired Fail rule makes the verdict FAIL,
// otherwise any fired Review rule makes it REVIEW, otherwise PASS. Each fired rule's
// Fact is stamped with the rule ID and collected as evidence (handed back to the
// agent on FAIL). It is deterministic: same (rules, ctx) ⇒ same Verdict.
//
// A defeasible weighting backend (toulmin h-Categoriser) would plug in here for
// genuine competing evidence (L2 AI consensus); the core path is level aggregation.
func Evaluate(rules []Rule, ctx Context) quest.Verdict {
	var facts []quest.Fact
	anyFail, anyReview := false, false
	failID, reviewID := "", ""
	for _, r := range rules {
		fired, fact := r.Check(ctx)
		if !fired {
			continue
		}
		fact.Rule = r.Meta.ID
		facts = append(facts, fact)
		if r.Meta.Level == LevelFail {
			anyFail = true
		}
		if r.Meta.Level == LevelFail && failID == "" {
			failID = r.Meta.ID
		}
		if r.Meta.Level == LevelReview {
			anyReview = true
		}
		if r.Meta.Level == LevelReview && reviewID == "" {
			reviewID = r.Meta.ID
		}
	}
	// RootCause = first fired Fail rule (or first Review rule when no Fail fired):
	// the decisive rule that produced the verdict, independent of Facts ordering.
	rootCause := failID
	if rootCause == "" {
		rootCause = reviewID
	}
	return verdictFromLevels(anyFail, anyReview, facts, rootCause)
}
