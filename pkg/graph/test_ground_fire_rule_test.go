//ff:func feature=graph type=helper control=sequence
//ff:what 테스트 헬퍼 — ctx.Grounds[name]이 비어있지 않을 때(즉 resolve된 본문이 위반을 증명할 때) 발동하는 tier-1 gate.Rule을 만든다. tier-1 규칙이 주입된 ground를 읽는지 확인하는 데 쓴다.

package graph

import (
	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// groundFireRule builds a tier-1 gate.Rule that fires when ctx.Grounds[name] is
// non-empty (i.e. the resolved body proves a violation). Used to confirm tier-1
// rules read the injected ground.
func groundFireRule(id, name string) gate.Rule {
	return gate.Rule{
		Meta: gate.RuleMeta{ID: id, Level: gate.LevelFail, Desc: id},
		Check: func(ctx gate.Context) (bool, quest.Fact) {
			return ctx.Grounds[name] != "", quest.Fact{Where: id, Actual: ctx.Grounds[name]}
		},
	}
}
