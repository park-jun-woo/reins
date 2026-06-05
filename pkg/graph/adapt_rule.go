//ff:func feature=graph type=helper control=sequence
//ff:what adaptRule — gate.Rule.Check(gate.Context)(bool, quest.Fact)를 toulmin func(Context, Specs)(bool, any)로 감싼다. toulmin ctx에서 gctxKey로 reins gate.Context를 꺼내 Check를 돌리고, 발동 시 (true, fact)를 — fact의 Rule ID를 ruleID로 스탬프해 — evidence로 반환한다. ctx 미주입 시 비발동.

package graph

import (
	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

// adaptRule wraps a gate.Rule's Check into a toulmin rule function. It retrieves
// the reins gate.Context from the toulmin context under gctxKey, runs Check, and
// on fire returns (true, fact) with fact.Rule stamped to ruleID as evidence. Each
// call to adaptRule returns a distinct closure (distinct fn value); uniqueness of
// the toulmin ruleID is guaranteed by the idSpec attached at registration.
func adaptRule(r gate.Rule, ruleID string) func(toulmin.Context, toulmin.Specs) (bool, any) {
	return func(ctx toulmin.Context, _ toulmin.Specs) (bool, any) {
		gv, ok := ctx.Get(gctxKey)
		if !ok {
			return false, nil
		}
		gctx, ok := gv.(gate.Context)
		if !ok {
			return false, nil
		}
		fired, fact := r.Check(gctx)
		if !fired {
			return false, nil
		}
		fact.Rule = ruleID
		return true, fact
	}
}
