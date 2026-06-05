//ff:func feature=graph type=helper control=sequence
//ff:what tier1GroundNames 빈 경로 증명 — 워런트는 skip하고 Needs가 없는 카운터도 skip하므로, ground 의존이 하나도 없으면 빈 슬라이스를 내는지 검증한다.

package graph

import (
	"testing"

	"github.com/park-jun-woo/reins/pkg/gate"
)

func TestTier1GroundNamesEmpty(t *testing.T) {
	g := NewGraph("g")
	pass := g.Warrant(alwaysTrue)                                           // warrant skipped
	g.Counter(fireRule("t0", gate.LevelFail), gate.LevelFail).Attacks(pass) // no Needs → skipped

	if names := g.tier1GroundNames(); len(names) != 0 {
		t.Fatalf("tier1GroundNames=%v want empty", names)
	}
}
