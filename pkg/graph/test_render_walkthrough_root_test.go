//ff:func feature=graph type=helper control=iteration dimension=1
//ff:what 공략집 렌더 증명 — fmt(Fail)·holder(Fail)·free(Review)가 발동하고 fmt가 holder·free를 Supersedes할 때 Evaluate가 FAIL을 내고 Verdict.Feedback에 ①결정타=fmt(상류 최상류)와 그 Fact, ②holder·free는 "fmt에 의해 superseded → 곁가지", ③"clear fmt" 다음 행동을 싣는지(문구 테이블) 검증한다.

package graph

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
)

func TestRenderWalkthroughRootCauseAndSuperseded(t *testing.T) {
	rFmt := fireRule("email-format", gate.LevelFail)
	rHolder := fireRule("source-lacks-email", gate.LevelFail)
	rFree := fireRule("freemail", gate.LevelReview)

	g := NewGraph("walk")
	pass := g.Warrant(alwaysTrue)
	fmtN := g.Counter(rFmt, gate.LevelFail).Attacks(pass)
	holder := g.Counter(rHolder, gate.LevelFail).Attacks(pass)
	free := g.Counter(rFree, gate.LevelReview).Attacks(pass)
	fmtN.Supersedes(holder, free)

	v := g.Evaluate(gate.Context{Submission: map[string]bool{
		"email-format": true, "source-lacks-email": true, "freemail": true,
	}})
	if v.Outcome != quest.OutFail {
		t.Fatalf("outcome=%s want FAIL", v.Outcome)
	}
	fb := v.Feedback
	if fb == "" {
		t.Fatal("Feedback empty, want walkthrough")
	}
	t.Logf("walkthrough:\n%s", fb)

	wants := []string{
		"FAIL. root cause = email-format (remaining active FAIL, upstream).",
		`Fact: where=email-format`,
		"source-lacks-email: superseded by email-format → side-branch.",
		"freemail: superseded by email-format → side-branch.",
		"→ to flip the verdict, clear email-format.",
	}
	for _, w := range wants {
		if !strings.Contains(fb, w) {
			t.Errorf("walkthrough missing %q\n--- got ---\n%s", w, fb)
		}
	}
}
