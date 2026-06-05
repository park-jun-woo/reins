//ff:func feature=graph type=helper control=sequence
//ff:what subGraph 포함/제외 증명 — include 술어로 선택된 카운터와 워런트만 살아남고, Attacks 엣지는 양끝이 살아남을 때만 재배선되며, supersession 엣지는 양끝이 포함될 때만 복사됨(상류 포함·하류 제외, 하류 포함·상류 제외 모두 미복사)을 검증한다. 결정적(네트워크 0).

package graph

import (
	"testing"

	"github.com/park-jun-woo/reins/pkg/gate"
)

func TestSubGraphIncludesWarrantsAndSelected(t *testing.T) {
	g := NewGraph("g")
	pass := g.Warrant(alwaysTrue)
	keep := g.Counter(fireRule("keep", gate.LevelFail), gate.LevelFail).Attacks(pass)
	drop := g.Counter(fireRule("drop", gate.LevelFail), gate.LevelFail).Attacks(pass)
	// supersession edges: keep→drop (down excluded), drop→keep (up excluded).
	keep.Supersedes(drop)
	drop.Supersedes(keep)

	// include only the "keep" counter; "drop" excluded. Warrant always included.
	sub := g.subGraph("sub", func(n *Node) bool { return n.shortName == keep.shortName })

	// Warrant + keep survive; drop does not.
	if nodeByID(sub, pass.id) == nil {
		t.Fatalf("subGraph: warrant not included")
	}
	if nodeByID(sub, keep.id) == nil {
		t.Fatalf("subGraph: kept counter not included")
	}
	if nodeByID(sub, drop.id) != nil {
		t.Fatalf("subGraph: dropped counter wrongly included")
	}
	if len(sub.nodes) != 2 {
		t.Fatalf("subGraph: node count=%d want 2 (warrant+keep)", len(sub.nodes))
	}

	// Attacks edge keep→pass survives (both included); drop→pass dropped.
	subKeep := nodeByID(sub, keep.id)
	if len(subKeep.attacks) != 1 || subKeep.attacks[0].id != pass.id {
		t.Fatalf("subGraph: keep->pass attack not rewired, attacks=%v", subKeep.attacks)
	}

	// Supersession: keep(up,incl)→drop(down,excl) NOT copied;
	// drop(up,excl)→keep NOT copied. So sub.supersedes must be empty.
	if len(sub.supersedes) != 0 {
		t.Fatalf("subGraph: supersedes should be empty when endpoints excluded, got %v", sub.supersedes)
	}
}
