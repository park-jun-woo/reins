//ff:func feature=graph type=helper control=iteration dimension=2
//ff:what tier1GroundNames — tier1 카운터(Needs 비어있지 않은 카운터)가 선언한 ground 이름들의 합집합을 등록 순서로(첫 등장 기준) 중복 없이 모은다. staged 평가가 tier0 통과 후 이 이름들만 lazy resolve해 ctx.Grounds에 채운다. 결정론(등록 순서).

package graph

// tier1GroundNames collects the union of ground names declared by tier-1 counters
// (counters whose Needs is non-empty), deduplicated and in registration order (by
// first appearance). After tier-0 passes, the staged evaluator resolves only these
// names lazily and populates ctx.Grounds. Deterministic (registration order).
func (g *Graph) tier1GroundNames() []string {
	seen := make(map[string]bool)
	var names []string
	for _, n := range g.nodes {
		if n.isWarrant || len(n.needs) == 0 {
			continue
		}
		for _, name := range n.needs {
			if seen[name] {
				continue
			}
			seen[name] = true
			names = append(names, name)
		}
	}
	return names
}
