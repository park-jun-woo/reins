//ff:func feature=graph type=helper control=iteration dimension=1
//ff:what 테스트 헬퍼 — 그래프 노드 목록을 훑어 주어진 id를 가진 노드를 반환한다. 없으면 nil. subGraph 결과 검증에 쓴다.

package graph

func nodeByID(g *Graph, id string) *Node {
	for _, n := range g.nodes {
		if n.id == id {
			return n
		}
	}
	return nil
}
