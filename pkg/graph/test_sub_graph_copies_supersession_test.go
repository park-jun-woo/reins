//ff:func feature=graph type=helper control=sequence
//ff:what subGraph supersession 복사 증명 — 상류·하류 노드가 모두 포함될 때만 supersession 엣지가 복사되고, 제외된 하류로의 엣지는 복사되지 않으며 set이 올바로 초기화됨을 검증한다(set-init 경로 포함). 결정적.

package graph

import (
	"testing"

	"github.com/park-jun-woo/reins/pkg/gate"
)

func TestSubGraphCopiesSupersessionWhenBothIncluded(t *testing.T) {
	g := NewGraph("g")
	pass := g.Warrant(alwaysTrue)
	up := g.Counter(fireRule("up", gate.LevelFail), gate.LevelFail).Attacks(pass)
	down := g.Counter(fireRule("down", gate.LevelReview), gate.LevelReview).Attacks(pass)
	extra := g.Counter(fireRule("extra", gate.LevelFail), gate.LevelFail).Attacks(pass)
	// up supersedes both down (included) and extra (excluded) — exercises the
	// "down excluded" continue inside an included-up branch and the set-init path.
	up.Supersedes(down, extra)

	include := func(n *Node) bool {
		return n.shortName == up.shortName || n.shortName == down.shortName
	}
	sub := g.subGraph("sub", include)

	set := sub.supersedes[up.id]
	if set == nil {
		t.Fatalf("subGraph: supersedes[up] not initialized")
	}
	if !set[down.id] {
		t.Fatalf("subGraph: up->down edge not copied")
	}
	if set[extra.id] {
		t.Fatalf("subGraph: up->extra wrongly copied (extra excluded)")
	}
	if len(set) != 1 {
		t.Fatalf("subGraph: supersedes[up] size=%d want 1", len(set))
	}
}
