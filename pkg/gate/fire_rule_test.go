//ff:func feature=gate type=helper control=sequence
//ff:what 테스트 헬퍼 — 주어진 ID·레벨·발동여부를 가진 Rule을 만든다(Check는 fire를 그대로 반환하고 Where=id인 Fact를 낸다).

package gate

import "github.com/park-jun-woo/reins/pkg/quest"

func fireRule(id string, lvl Level, fire bool) Rule {
	return Rule{
		Meta: RuleMeta{ID: id, Level: lvl, Desc: id},
		Check: func(Context) (bool, quest.Fact) {
			return fire, quest.Fact{Where: id}
		},
	}
}
