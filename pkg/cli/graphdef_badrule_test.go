//ff:func feature=cli type=helper control=sequence
//ff:what graphDef.badRule. 제출물 텍스트가 "bad"면 발동(FAIL)하는 공유 카운터를 돌려준다(테스트 더블).

package cli

import (
	"strings"

	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// badRule is the shared counter: it fires (FAIL) when the submission text is "bad".
func (graphDef) badRule() gate.Rule {
	return gate.Rule{
		Meta: gate.RuleMeta{ID: "not-bad", Level: gate.LevelFail, Desc: "submission must not be bad"},
		Check: func(ctx gate.Context) (bool, quest.Fact) {
			if s, _ := ctx.Submission.(string); strings.TrimSpace(s) == "bad" {
				return true, quest.Fact{Where: "body", Expected: "good", Actual: "bad"}
			}
			return false, quest.Fact{}
		},
	}
}
