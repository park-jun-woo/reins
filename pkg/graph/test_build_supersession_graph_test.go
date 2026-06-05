//ff:func feature=graph type=helper control=sequence
//ff:what buildSupersessionGraph 테스트헬퍼 — fmt가 holder+free를, free가 holder를 Supersedes하는 3노드 그래프와 id→Node 맵을 만든다(applySupersession 단위증명 고정 입력).

package graph

func buildSupersessionGraph() (*Graph, map[string]*Node) {
	g := NewGraph("sup")
	fmtN := mkNode(g, "fmt")
	holder := mkNode(g, "holder")
	free := mkNode(g, "free")
	g.supersedes["fmt"] = map[string]bool{"holder": true, "free": true}
	g.supersedes["free"] = map[string]bool{"holder": true}
	g.nodes = []*Node{fmtN, holder, free}
	return g, map[string]*Node{"fmt": fmtN, "holder": holder, "free": free}
}
