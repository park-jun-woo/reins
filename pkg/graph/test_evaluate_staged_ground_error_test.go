//ff:func feature=graph type=helper control=sequence
//ff:what staged 평가의 ground resolve 실패 경로 — tier0 통과 후 provider가 에러를 내면(접근 불가) terminal FAIL Verdict로 환원되고 그 Fact(Rule="ground", Where=ground 이름, Actual=에러)를 담는지 검증한다(열린결정 #8). resolve 에러는 결정론 FAIL Fact.

package graph

import (
	"errors"
	"testing"

	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/ground"
	"github.com/park-jun-woo/reins/pkg/quest"
)

func TestEvaluateStagedGroundError(t *testing.T) {
	g := NewGraph("staged-err")
	pass := g.Warrant(alwaysTrue)
	g.Counter(fireRule("email-format", gate.LevelFail), gate.LevelFail).Attacks(pass)
	g.Counter(groundFireRule("source-lacks-email", "source-body"), gate.LevelFail).
		Attacks(pass).Needs("source-body")

	fr := &stagedFakeResolver{body: "x"}
	snap := ground.NewSnapshot(fr)

	// Provider that always errors (e.g. source unreachable).
	provider := func(name string, _ gate.Context, _ *ground.Snapshot) (string, error) {
		return "", errors.New("unreachable")
	}

	// tier 0 clean → ground resolve attempted → error → terminal FAIL.
	v := g.EvaluateStaged(gate.Context{Submission: map[string]bool{}}, snap, provider)
	if v.Outcome != quest.OutFail {
		t.Fatalf("ground error: outcome=%s want FAIL", v.Outcome)
	}
	if len(v.Facts) != 1 || v.Facts[0].Rule != "ground" || v.Facts[0].Where != "source-body" {
		t.Fatalf("ground error: want one ground/source-body fact, got %v", v.Facts)
	}
}
