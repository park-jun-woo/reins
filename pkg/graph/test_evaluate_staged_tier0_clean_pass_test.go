//ff:func feature=graph type=helper control=sequence
//ff:what staged 평가 PASS 경로 — clean tier-0 + 위반을 일으키지 않는 tier-1 ground(빈 본문) → PASS이되 ground는 여전히 1회 resolve됨을 증명한다.

package graph

import (
	"testing"

	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/ground"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// TestEvaluateStagedTier0CleanPass — clean tier 0 + tier-1 grounds that do not
// trigger a violation → PASS, but the ground was still resolved once.
func TestEvaluateStagedTier0CleanPass(t *testing.T) {
	g := NewGraph("staged-pass")
	pass := g.Warrant(alwaysTrue)
	g.Counter(fireRule("email-format", gate.LevelFail), gate.LevelFail).Attacks(pass)
	// tier-1 rule that fires only if ground == "BAD"; our body is empty → no fire.
	g.Counter(groundFireRule("source-lacks-email", "source-body"), gate.LevelFail).
		Attacks(pass).Needs("source-body")

	fr := &stagedFakeResolver{body: ""} // empty body → tier-1 rule does not fire
	snap := ground.NewSnapshot(fr)

	v := g.EvaluateStaged(gate.Context{Submission: map[string]bool{}}, snap, stagedProvider)
	if fr.fetchCalls != 1 {
		t.Fatalf("clean pass: Fetch called %d times, want 1", fr.fetchCalls)
	}
	if v.Outcome != quest.OutPass {
		t.Fatalf("clean pass: outcome=%s want PASS", v.Outcome)
	}
}
