//ff:func feature=graph type=helper control=iteration dimension=2
//ff:what supersedesReach — 노드 id가 supersession 그래프에서 (직·간접) 흡수하는 하류 노드 수를 센다. 상류일수록 reach가 크다 → 결정타(위상 최상류) 선택의 순위 키. 자기 자신은 제외, 사이클은 visited로 차단(toulmin acyclic 전제와 무관히 안전).

package graph

// supersedesReach counts how many downstream nodes the given node id absorbs,
// directly or transitively, in the supersession graph. A more upstream node has a
// larger reach — this is the ranking key for selecting the root cause (the
// topologically uppermost). The node itself is excluded; cycles are guarded by the
// visited set.
func (g *Graph) supersedesReach(id string) int {
	visited := make(map[string]bool)
	var walk func(string)
	walk = func(cur string) {
		for down := range g.supersedes[cur] {
			if visited[down] {
				continue
			}
			visited[down] = true
			walk(down)
		}
	}
	walk(id)
	return len(visited)
}
