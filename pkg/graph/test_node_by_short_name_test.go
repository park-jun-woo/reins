//ff:func feature=graph type=helper control=sequence
//ff:what nodeByShortName 증명 — shortName이 일치하는 노드를 되찾고, 미일치 시 nil을 반환하는지 검증한다.

package graph

import "testing"

func TestNodeByShortName(t *testing.T) {
	g := NewGraph("byname")
	a := &Node{graph: g, id: "a", shortName: "alpha#0"}
	b := &Node{graph: g, id: "b", shortName: "beta#1"}
	g.nodes = []*Node{a, b}

	if got := g.nodeByShortName("alpha#0"); got != a {
		t.Fatalf("alpha#0: got %v want node a", got)
	}
	if got := g.nodeByShortName("beta#1"); got != b {
		t.Fatalf("beta#1: got %v want node b", got)
	}
	if got := g.nodeByShortName("missing"); got != nil {
		t.Fatalf("missing: got %v want nil", got)
	}
}
