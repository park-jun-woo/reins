//ff:func feature=graph type=helper control=sequence
//ff:what REVIEW 공략집 증명 — 잔존이 Review-Level 카운터뿐이면 Evaluate가 REVIEW를 내고 Verdict.Feedback의 결정타가 그 Review 카운터(focal point)로 렌더되는지 검증한다(잔존 Fail이 없을 때 첫 잔존을 결정타로).

package graph

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
)

func TestRenderWalkthroughReviewFocus(t *testing.T) {
	rFree := fireRule("freemail", gate.LevelReview)
	g := NewGraph("review")
	pass := g.Warrant(alwaysTrue)
	g.Counter(rFree, gate.LevelReview).Attacks(pass)

	v := g.Evaluate(gate.Context{Submission: map[string]bool{"freemail": true}})
	if v.Outcome != quest.OutReview {
		t.Fatalf("outcome=%s want REVIEW", v.Outcome)
	}
	if !strings.Contains(v.Feedback, "REVIEW. root cause = freemail") {
		t.Fatalf("review walkthrough = %q", v.Feedback)
	}
}
