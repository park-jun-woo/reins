//ff:func feature=cli type=helper control=sequence
//ff:what graphDef.Rules. badRule 카운터 하나로 된 카탈로그를 돌려준다(rules 명령·후방호환 감사용 테스트 더블).

package cli

import "github.com/park-jun-woo/reins/pkg/gate"

func (d graphDef) Rules() []gate.Rule { return []gate.Rule{d.badRule()} }
