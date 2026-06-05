//ff:type feature=gate type=model
//ff:what 위반 탐지기. Check가 fired=true + Fact를 내면 문제 발견(안 터지면 문제 없음). 발동 규칙들의 레벨 집계가 verdict를 결정한다(Evaluate 참조).

package gate

import "github.com/park-jun-woo/reins/pkg/quest"

// Rule is a violation detector. Check returns fired=true with a Fact when it finds a
// problem; no fire means no problem. The aggregate of fired rules' levels decides
// the verdict (see Evaluate).
type Rule struct {
	Meta  RuleMeta
	Check func(ctx Context) (fired bool, fact quest.Fact)
}
