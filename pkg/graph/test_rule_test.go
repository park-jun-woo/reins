//ff:func feature=graph type=helper control=sequence
//ff:what 테스트 헬퍼 — ID·레벨을 가진 gate.Rule을 만든다. Check는 ctx.Submission(map[string]bool)에서 자기 ID를 찾아 발동여부를 정하고 Where=id인 Fact를 낸다. 같은 ctx로 두 평가 경로(graph vs gate)를 비교하기 위한 결정적 규칙.

package graph

import (
	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// fireRule builds a gate.Rule with the given ID and Level whose Check fires when
// ctx.Submission (a map[string]bool) marks its ID true, emitting a Fact with
// Where=id. Deterministic — used to compare the graph and gate evaluation paths.
func fireRule(id string, lvl gate.Level) gate.Rule {
	return gate.Rule{
		Meta: gate.RuleMeta{ID: id, Level: lvl, Desc: id},
		Check: func(ctx gate.Context) (bool, quest.Fact) {
			m, _ := ctx.Submission.(map[string]bool)
			return m[id], quest.Fact{Where: id}
		},
	}
}
