//ff:func feature=graph type=helper control=sequence
//ff:what REVIEW 잔존 병기 증명 — FAIL root(rule-fail)와 함께 supersession 없이 잔존한 REVIEW 카운터(rule-review)가 Facts에만 있고 공략집에서 빠지지 않도록, walkthrough에 "remaining side-branch"로 병기되는지 검증한다(피드백 패리티).

package graph

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// TestRenderWalkthroughReviewSibling: a REVIEW counter that remains active next to
// a FAIL root (no supersession) appears in the walkthrough as a remaining
// side-branch — feedback parity with Facts, which already list it.
func TestRenderWalkthroughReviewSibling(t *testing.T) {
	rFail := fireRule("rule-fail", gate.LevelFail)
	rReview := fireRule("rule-review", gate.LevelReview)

	g := NewGraph("review-sibling")
	pass := g.Warrant(alwaysTrue)
	g.Counter(rFail, gate.LevelFail).Attacks(pass)
	g.Counter(rReview, gate.LevelReview).Attacks(pass)

	v := g.Evaluate(gate.Context{Submission: map[string]bool{"rule-fail": true, "rule-review": true}})
	if v.Outcome != quest.OutFail {
		t.Fatalf("outcome=%s want FAIL", v.Outcome)
	}
	fb := v.Feedback
	t.Logf("walkthrough:\n%s", fb)
	if !strings.Contains(fb, "root cause = rule-fail") {
		t.Errorf("want rule-fail as root cause, got:\n%s", fb)
	}
	if !strings.Contains(fb, "rule-review: remaining side-branch") {
		t.Errorf("want surviving REVIEW counter rule-review as remaining side-branch, got:\n%s", fb)
	}
}
