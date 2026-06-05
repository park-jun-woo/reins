//ff:func feature=graph type=helper control=sequence
//ff:what resolveGrounds ① — tier1 ground 이름이 0개면 (nil,nil)을 내고 provider를 한 번도 호출하지 않음(네트워크 0)을 호출 카운터로 증명한다.

package graph

import (
	"testing"

	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/ground"
)

func TestResolveGroundsNoNames(t *testing.T) {
	g := NewGraph("g")
	pass := g.Warrant(alwaysTrue)
	g.Counter(fireRule("t0", gate.LevelFail), gate.LevelFail).Attacks(pass) // no Needs

	fr := &stagedFakeResolver{body: "x"}
	snap := ground.NewSnapshot(fr)
	calls := 0
	provider := func(name string, _ gate.Context, _ *ground.Snapshot) (string, error) {
		calls++
		return "v", nil
	}

	grounds, v := g.resolveGrounds(gate.Context{}, snap, provider)
	if grounds != nil || v != nil {
		t.Fatalf("no names: got (%v,%v) want (nil,nil)", grounds, v)
	}
	if calls != 0 {
		t.Fatalf("no names: provider called %d times, want 0", calls)
	}
}
