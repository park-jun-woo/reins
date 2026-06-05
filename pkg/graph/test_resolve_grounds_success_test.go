//ff:func feature=graph type=helper control=sequence
//ff:what resolveGrounds ② — ground 이름이 있고 provider가 성공하면 grounds 맵이 채워지고 Verdict는 nil임을 증명한다.

package graph

import (
	"testing"

	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/ground"
)

func TestResolveGroundsSuccess(t *testing.T) {
	g := NewGraph("g")
	pass := g.Warrant(alwaysTrue)
	g.Counter(groundFireRule("c", "body"), gate.LevelFail).Attacks(pass).Needs("body")

	fr := &stagedFakeResolver{body: "x"}
	snap := ground.NewSnapshot(fr)
	provider := func(name string, _ gate.Context, _ *ground.Snapshot) (string, error) {
		return "resolved-" + name, nil
	}

	grounds, v := g.resolveGrounds(gate.Context{}, snap, provider)
	if v != nil {
		t.Fatalf("success: verdict=%v want nil", v)
	}
	if grounds["body"] != "resolved-body" {
		t.Fatalf("success: grounds=%v want body=resolved-body", grounds)
	}
}
