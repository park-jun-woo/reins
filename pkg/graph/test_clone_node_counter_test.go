//ff:func feature=graph type=helper control=sequence
//ff:what cloneNode 카운터 경로 증명 — 카운터 노드를 dst 그래프에 복제하면 id·shortName에 더해 level·needs까지 복사되고 graph/tr이 dst로 바인딩되며 dst.nodes에 추가되는지 검증한다(네트워크 0, 결정적).

package graph

import (
	"testing"

	"github.com/park-jun-woo/reins/pkg/gate"
)

func TestCloneNodeCounter(t *testing.T) {
	src := NewGraph("src")
	src.Warrant(alwaysTrue)
	c := src.Counter(fireRule("c1", gate.LevelReview), gate.LevelReview).Needs("g1", "g2")

	dst := NewGraph("dst")
	got := dst.cloneNode(c)

	if got.isWarrant {
		t.Fatalf("cloned counter: isWarrant=true want false")
	}
	if got.id != c.id || got.shortName != c.shortName {
		t.Fatalf("cloned counter: id=%q short=%q want id=%q short=%q",
			got.id, got.shortName, c.id, c.shortName)
	}
	if got.level != c.level {
		t.Fatalf("cloned counter: level=%v want %v", got.level, c.level)
	}
	if len(got.needs) != 2 || got.needs[0] != "g1" || got.needs[1] != "g2" {
		t.Fatalf("cloned counter: needs=%v want [g1 g2]", got.needs)
	}
	if got.graph != dst || got.tr == nil {
		t.Fatalf("cloned counter: graph/tr not bound to dst")
	}
	if len(dst.nodes) != 1 || dst.nodes[0] != got {
		t.Fatalf("cloned counter: dst.nodes=%v want [got]", dst.nodes)
	}
}
