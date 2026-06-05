//ff:func feature=graph type=helper control=sequence
//ff:what Evaluate REVIEW 분기증명 — 잔존 카운터가 전부 Review면 REVIEW(+Facts)로 환원하는지 검증한다(기존 basic/supersession이 PASS·FAIL 분기를 덮음).

package graph

import (
	"testing"

	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
)

func TestEvaluateAllReview(t *testing.T) {
	rFree := fireRule("free", gate.LevelReview)
	rGroup := fireRule("group", gate.LevelReview)

	g := NewGraph("review")
	pass := g.Warrant(alwaysTrue)
	g.Counter(rFree, gate.LevelReview).Attacks(pass)
	g.Counter(rGroup, gate.LevelReview).Attacks(pass)

	v := g.Evaluate(gate.Context{Submission: map[string]bool{"free": true, "group": true}})
	if v.Outcome != quest.OutReview {
		t.Fatalf("outcome=%s want REVIEW", v.Outcome)
	}
	if len(v.Facts) != 2 {
		t.Fatalf("expected 2 facts, got %d (%+v)", len(v.Facts), v.Facts)
	}
}
