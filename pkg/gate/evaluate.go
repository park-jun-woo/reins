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
	for _, r := range rules {
		fired, fact := r.Check(ctx)
		if !fired {
			continue
		}
		fact.Rule = r.Meta.ID
		facts = append(facts, fact)
		switch r.Meta.Level {
		case LevelFail:
			anyFail = true
		case LevelReview:
			anyReview = true
		}
	}
	switch {
	case anyFail:
		return quest.Verdict{Outcome: quest.OutFail, Facts: facts}
	case anyReview:
		return quest.Verdict{Outcome: quest.OutReview, Facts: facts}
	default:
		return quest.Verdict{Outcome: quest.OutPass}
	}
}

// Catalog returns the metas of the given rules — the auto-generated rulebook that
// the cli `rules` command prints (every cheese this gate blocks, by ID and level).
func Catalog(rules []Rule) []RuleMeta {
	metas := make([]RuleMeta, 0, len(rules))
	for _, r := range rules {
		metas = append(metas, r.Meta)
	}
	return metas
}
