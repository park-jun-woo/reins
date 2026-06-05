//ff:func feature=cli type=helper control=sequence
//ff:what graphDef.Evaluate. defeat 그래프에서 verdict를 읽는다(Evaluator 경로). graph.FromRules가 edge-zero 등가 그래프(FAIL 카운터에 공격받는 tautology PASS warrant 하나)를 만들어 submit의 Evaluator 분기가 실제 pkg/graph로 작동함을 검증한다(테스트 더블).

package cli

import (
	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/graph"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// Evaluate reads the verdict from a defeat graph (Evaluator path). graph.FromRules
// builds the edge-zero equivalent (one tautology PASS warrant attacked by the FAIL
// counter), so submit's Evaluator branch is exercised against the real pkg/graph.
func (d graphDef) Evaluate(ctx gate.Context) quest.Verdict {
	return graph.FromRules(d.Rules()).Evaluate(ctx)
}
