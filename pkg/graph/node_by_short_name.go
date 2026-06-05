//ff:func feature=graph type=helper control=iteration dimension=1
//ff:what nodeByShortName — 수집된 trace 엔트리(Name=short)를 노드로 되찾는다. idSpec가 short name을 노드별로 고유하게 만들므로 매칭은 일의적이다. 미발견 시 nil.

package graph

// nodeByShortName returns the node whose shortName matches name, or nil. The
// idSpec makes each node's short name unique, so the match is unambiguous.
func (g *Graph) nodeByShortName(name string) *Node {
	for _, n := range g.nodes {
		if n.shortName == name {
			return n
		}
	}
	return nil
}
