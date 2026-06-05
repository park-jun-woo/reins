//ff:func feature=cli type=helper control=sequence
//ff:what errDef.Rules. 규칙이 없는 게이트를 흉내내 nil을 돌려준다(테스트 더블).

package cli

import "github.com/park-jun-woo/reins/pkg/gate"

func (d errDef) Rules() []gate.Rule { return nil }
