//ff:func feature=graph type=helper control=iteration dimension=1
//ff:what renderWalkthrough 직접 단위증명 — 결정타 Fact·곁가지(잔존 Fail)·superseded 이벤트·다음 행동 줄을 한 입력으로 모두 렌더하는지, 그리고 잔존 0이면 빈 문자열을 내는지(PASS 분기) 검증한다. Evaluate를 거치지 않고 함수 직접 호출.

package graph

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
)

func TestRenderWalkthroughUnit(t *testing.T) {
	g, byID := buildSupersessionGraph()
	byID["fmt"].level = gate.LevelFail
	byID["holder"].level = gate.LevelFail
	byID["free"].level = gate.LevelReview

	// empty remaining -> "".
	if s := g.renderWalkthrough("FAIL", nil, nil); s != "" {
		t.Fatalf("empty remaining = %q, want \"\"", s)
	}

	root := activeCounter{node: byID["fmt"], fact: quest.Fact{Rule: "email-format", Where: "email", Expected: "user@domain", Actual: "not-an-email"}}
	sib := activeCounter{node: byID["holder"], fact: quest.Fact{Rule: "source-lacks-email", Where: "body"}}
	// remaining holds two Fails: fmt (upstream) is root, holder is a remaining side-branch.
	remaining := []activeCounter{root, sib}
	// free was absorbed by fmt.
	events := []supersessionEvent{{
		absorbed: activeCounter{node: byID["free"], fact: quest.Fact{Rule: "freemail"}},
		by:       activeCounter{node: byID["fmt"], fact: quest.Fact{Rule: "email-format"}},
	}}

	out := g.renderWalkthrough("FAIL", remaining, events)
	t.Logf("walkthrough:\n%s", out)
	for _, w := range []string{
		"FAIL. root cause = email-format (remaining active FAIL, upstream).",
		`Fact: where=email expected="user@domain" actual="not-an-email"`,
		"source-lacks-email: remaining side-branch (where=body",
		"freemail: superseded by email-format → side-branch.",
		"→ to flip the verdict, clear email-format.",
	} {
		if !strings.Contains(out, w) {
			t.Errorf("missing %q in:\n%s", w, out)
		}
	}
}
