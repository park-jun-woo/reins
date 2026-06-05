//ff:func feature=graph type=helper control=sequence
//ff:what mkNode 테스트헬퍼 — 그래프 g에 속한 id/shortName=id 노드 1개를 만들어 반환한다(applySupersession 단위증명용).

package graph

func mkNode(g *Graph, id string) *Node {
	return &Node{graph: g, id: id, shortName: id}
}
