//ff:func feature=graph type=helper control=sequence
//ff:what Counter — 위반 카운터 노드를 그래프에 추가한다. gate.Rule.Check를 어댑터로 감싸 toulmin g.Counter(adapted).With(고유 idSpec)로 등록하고, reins Level을 노드 메타로 보유한다(toulmin Strength 미사용). counter→warrant 공격 배선은 호출측이 Node.Attacks로 한다. 어댑터의 fact는 gate 규칙 ID로 스탬프된다.

package graph

import (
	"strconv"

	"github.com/park-jun-woo/reins/pkg/gate"
)

// Counter adds a violation counter node to the graph. It wraps rule.Check via an
// adapter, registers it with toulmin (g.Counter(adapted).With(unique idSpec)), and
// holds the reins Level as node meta (toulmin Strength is unused). Wiring the
// counter→warrant attack edge is the caller's job via Node.Attacks. The adapter's
// Fact is stamped with the gate rule's ID.
func (g *Graph) Counter(rule gate.Rule, level gate.Level) *Node {
	idStr := rule.Meta.ID + "#n" + strconv.Itoa(len(g.nodes))
	spec := idSpec{ID: idStr}
	adapted := adaptRule(rule, rule.Meta.ID)
	id := ruleIDFor(adapted, spec)
	tr := g.tg.Counter(adapted).With(spec)
	n := &Node{
		graph:     g,
		tr:        tr,
		id:        id,
		shortName: shortNameFor(id),
		level:     level,
		isWarrant: false,
		fn:        adapted,
		spec:      spec,
	}
	g.nodes = append(g.nodes, n)
	return n
}
