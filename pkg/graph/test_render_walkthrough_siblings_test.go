//ff:func feature=graph type=helper control=sequence
//ff:what 곁가지 Fail 병기 증명 — supersession 관계 없는 잔존 Fail이 둘이면 등록순 첫째(rule-a)가 결정타가 되고 나머지(rule-b)는 "remaining side-branch"로 병기되는지 검증한다(결정타 단일화·등록순 폴백).

package graph

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
)

func TestRenderWalkthroughSiblingFails(t *testing.T) {
	rA := fireRule("rule-a", gate.LevelFail)
	rB := fireRule("rule-b", gate.LevelFail)

	g := NewGraph("siblings")
	pass := g.Warrant(alwaysTrue)
	g.Counter(rA, gate.LevelFail).Attacks(pass)
	g.Counter(rB, gate.LevelFail).Attacks(pass)

	v := g.Evaluate(gate.Context{Submission: map[string]bool{"rule-a": true, "rule-b": true}})
	if v.Outcome != quest.OutFail {
		t.Fatalf("outcome=%s want FAIL", v.Outcome)
	}
	fb := v.Feedback
	t.Logf("walkthrough:\n%s", fb)
	if !strings.Contains(fb, "root cause = rule-a") {
		t.Errorf("want rule-a as root cause, got:\n%s", fb)
	}
	if !strings.Contains(fb, "rule-b: remaining side-branch") {
		t.Errorf("want rule-b as remaining side-branch, got:\n%s", fb)
	}
}
