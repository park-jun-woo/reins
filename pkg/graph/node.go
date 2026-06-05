//ff:type feature=graph type=model
//ff:what Node — Graph 내 노드의 reins-side 핸들. toulmin *Rule(엣지 배선용)·reins Level 메타·고유 ruleID·trace 매칭용 short name·Fact 운반 여부(워런트인지)를 보유한다. staged 평가를 위해 등록 fn·idSpec·공격 대상 노드·Needs(ground 의존 선언)도 보관 — 이들로 tier 서브그래프를 재구성한다. Attacks/Supersedes/Needs 체이닝을 지원한다.

package graph

import (
	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

// Node is a reins-side handle to a node registered in a Graph. It carries the
// toulmin *Rule (for wiring defeat edges), the reins Level meta (the basis of the
// crisp verdict — toulmin Strength is unused), the node's unique toulmin ruleID, its
// short name (for trace matching), and whether it is a warrant. For staged
// evaluation it also retains the registration fn, idSpec, attack target nodes, and
// the declared ground deps (Needs) so a tier sub-graph can be rebuilt from these.
type Node struct {
	graph     *Graph
	tr        *toulmin.Rule
	id        string     // toulmin ruleID (funcID#spec)
	shortName string     // shortName(id) — matches collected trace entry Name
	level     gate.Level // reins Level meta (Fail/Review); meaningful for counters
	isWarrant bool

	fn      func(toulmin.Context, toulmin.Specs) (bool, any) // registration fn (for staged rebuild)
	spec    idSpec                                           // registration idSpec
	attacks []*Node                                          // nodes this counter attacks (counter→warrant)
	needs   []string                                         // declared ground deps; non-empty ⇒ tier-1
}
