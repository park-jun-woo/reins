//ff:func feature=graph type=helper control=sequence
//ff:what resolveGrounds ③ — provider가 에러를 내면 grounds는 nil이고 terminal FAIL Verdict(Rule="ground", Where=이름, Actual=에러)를 내는지 증명한다.

package graph

import (
	"errors"
	"testing"

	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/ground"
	"github.com/park-jun-woo/reins/pkg/quest"
)

func TestResolveGroundsProviderError(t *testing.T) {
	g := NewGraph("g")
	pass := g.Warrant(alwaysTrue)
	g.Counter(groundFireRule("c", "body"), gate.LevelFail).Attacks(pass).Needs("body")

	fr := &stagedFakeResolver{body: "x"}
	snap := ground.NewSnapshot(fr)
	provider := func(name string, _ gate.Context, _ *ground.Snapshot) (string, error) {
		return "", errors.New("unreachable")
	}

	grounds, v := g.resolveGrounds(gate.Context{}, snap, provider)
	if grounds != nil {
		t.Fatalf("error: grounds=%v want nil", grounds)
	}
	if v == nil || v.Outcome != quest.OutFail {
		t.Fatalf("error: verdict=%v want FAIL", v)
	}
	if len(v.Facts) != 1 || v.Facts[0].Rule != "ground" ||
		v.Facts[0].Where != "body" || v.Facts[0].Actual != "unreachable" {
		t.Fatalf("error: facts=%v want one ground/body/unreachable fact", v.Facts)
	}
}
