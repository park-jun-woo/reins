//ff:func feature=graph type=helper control=sequence
//ff:what alwaysTrue — 항상 활성인 워런트 함수 헬퍼. tautology PASS 워런트의 fn으로 쓴다(위반 카운터들이 이걸 공격). toulmin isWarrant가 attacker를 결과에서 제외하므로, 항상-활성 워런트가 있어야 Evaluate가 결과 1개를 내고 그 trace에 카운터별 Activated가 실린다.

package graph

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// alwaysTrue is an always-active warrant function. Use it as the fn of the
// tautology PASS warrant that violation counters attack. toulmin's isWarrant
// excludes attacker nodes from results, so an always-active warrant is required
// for Evaluate to emit one result whose trace carries each counter's Activated.
func alwaysTrue(_ toulmin.Context, _ toulmin.Specs) (bool, any) {
	return true, nil
}
