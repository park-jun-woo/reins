//ff:func feature=graph type=helper control=selection
//ff:what cloneNode — 원 노드(src)를 대상 그래프(dst)의 새 toulmin 그래프에 같은 fn·idSpec으로 재등록한다. 같은 fn+spec이라 toulmin ruleID/shortName이 원본과 동일 → supersession 맵(ID 키)·trace 매칭이 그대로 옮겨진다. 워런트는 tg.Rule, 카운터는 tg.Counter로 등록하고 level/needs를 복사한다. 엣지(Attacks)는 호출측이 별도 재배선한다. staged tier 서브그래프 재구성의 단위.

package graph

// cloneNode re-registers src into dst's fresh toulmin graph with the same fn and
// idSpec. Because the fn+spec are identical, the toulmin ruleID (and shortName)
// match the original, so the supersession map (keyed by ID) and trace matching carry
// over unchanged. A warrant is registered via tg.Rule, a counter via tg.Counter;
// level and needs are copied. Attack edges are rewired separately by the caller. It
// is the unit of staged-tier sub-graph reconstruction.
func (dst *Graph) cloneNode(src *Node) *Node {
	var n *Node
	if src.isWarrant {
		tr := dst.tg.Rule(src.fn).With(src.spec)
		n = &Node{
			graph:     dst,
			tr:        tr,
			id:        src.id,
			shortName: src.shortName,
			isWarrant: true,
			fn:        src.fn,
			spec:      src.spec,
		}
	} else {
		tr := dst.tg.Counter(src.fn).With(src.spec)
		n = &Node{
			graph:     dst,
			tr:        tr,
			id:        src.id,
			shortName: src.shortName,
			level:     src.level,
			isWarrant: false,
			fn:        src.fn,
			spec:      src.spec,
			needs:     src.needs,
		}
	}
	dst.nodes = append(dst.nodes, n)
	return n
}
