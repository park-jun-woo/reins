//ff:func feature=cli type=helper control=sequence
//ff:what stubDef.Rules. 제출물이 "bad"면 발동하는 단일 FAIL 규칙 카탈로그를 돌려준다(게이트 평가 테스트 더블).

package cli

import (
	"strings"

	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
)

func (stubDef) Rules() []gate.Rule {
	return []gate.Rule{{
		Meta: gate.RuleMeta{ID: "not-bad", Level: gate.LevelFail, Desc: "submission must not be bad"},
		Check: func(ctx gate.Context) (bool, quest.Fact) {
			if s, _ := ctx.Submission.(string); strings.TrimSpace(s) == "bad" {
				return true, quest.Fact{Where: "body", Expected: "good", Actual: "bad"}
			}
			return false, quest.Fact{}
		},
	}}
}
