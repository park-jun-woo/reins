//ff:type feature=gate type=model
//ff:what 선택적 Definition 확장. def가 Evaluator이면 게이트는 Rules() 카탈로그 대신 자체 Evaluate(defeat 그래프 등)로 Verdict를 계산한다. gate는 graph/toulmin을 import하지 않는다(toulmin-free 유지) — 메서드 시그니처만 선언. cli가 def.(Evaluator) 타입단언으로 분기.

package gate

import "github.com/park-jun-woo/reins/pkg/quest"

// Evaluator is an optional Definition extension: a gate that computes its own
// Verdict (e.g. via a defeat graph) instead of the Rules() catalog. pkg/gate does
// not import pkg/graph or toulmin — it only declares the method, and the cli wiring
// branches on a def.(Evaluator) assertion. Rules() is still implemented for the
// `rules` command (the audit catalog); only Evaluate uses the graph.
type Evaluator interface {
	Evaluate(ctx Context) quest.Verdict
}
