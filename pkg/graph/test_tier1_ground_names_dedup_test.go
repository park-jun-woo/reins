//ff:func feature=graph type=helper control=sequence
//ff:what tier1GroundNames dedup·순서 증명 — 여러 카운터의 Needs를 모으되 중복 ground 이름은 한 번만 담고 첫 등장 순서를 보존하는지 검증한다(seen continue 경로).

package graph

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/reins/pkg/gate"
)

func TestTier1GroundNamesDedupAndOrder(t *testing.T) {
	g := NewGraph("g")
	pass := g.Warrant(alwaysTrue)
	g.Counter(fireRule("t0", gate.LevelFail), gate.LevelFail).Attacks(pass) // no Needs
	g.Counter(fireRule("a", gate.LevelFail), gate.LevelFail).
		Attacks(pass).Needs("body", "mx") // first appearance: body, mx
	g.Counter(fireRule("b", gate.LevelFail), gate.LevelFail).
		Attacks(pass).Needs("mx", "title") // mx duplicate (seen) → skipped; title new

	got := g.tier1GroundNames()
	want := []string{"body", "mx", "title"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("tier1GroundNames=%v want %v", got, want)
	}
}
