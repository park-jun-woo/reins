//ff:func feature=graph type=helper control=iteration dimension=2
//ff:what renderWalkthrough — 평가 중 수집한 잔존 카운터·superseded 이벤트로 에이전트 직통 "공략집"을 만든다. ①결정타(잔존 활성 최상류 Fail)와 그 결정론 Fact, ②다른 잔존 Fail은 "+곁가지(잔존)"로 병기, ③superseded된 활성 카운터는 "X에 의해 superseded → 곁가지" 한 줄(완전 숨김 아님), ④"판정을 뒤집으려면 <결정타>를 끄라" 다음 행동. 잎은 전부 규칙이 낸 quest.Fact. remaining 없으면(PASS) 빈 문자열.

package graph

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/reins/pkg/quest"
)

// renderWalkthrough builds the agent-facing strategy text ("walkthrough") from the
// remaining counters and the superseded events gathered during evaluation:
// (1) the root cause (the uppermost remaining active Fail) with its deterministic
// Fact, (2) any other remaining Fail listed as a remaining side-branch, (3) each
// superseded active counter as a one-line "superseded by X → side-branch" (not fully
// hidden), and (4) a "to flip the verdict, clear <root cause>" next action. Every
// leaf is a rule-emitted quest.Fact. Returns "" when nothing remains (PASS).
func (g *Graph) renderWalkthrough(outcome string, remaining []activeCounter, events []supersessionEvent) string {
	if len(remaining) == 0 {
		return ""
	}
	root, siblings := g.selectRootCause(remaining)

	var b strings.Builder
	fmt.Fprintf(&b, "%s. root cause = %s (remaining active %s, upstream).\n",
		outcome, counterName(root), root.node.level)
	if f := root.fact; f != (quest.Fact{}) {
		fmt.Fprintf(&b, "  Fact: where=%s expected=%q actual=%q\n", f.Where, f.Expected, f.Actual)
	}
	for _, sib := range siblings {
		f := sib.fact
		fmt.Fprintf(&b, "  + %s: remaining side-branch (where=%s expected=%q actual=%q).\n",
			counterName(sib), f.Where, f.Expected, f.Actual)
	}
	for _, ev := range events {
		fmt.Fprintf(&b, "  %s: superseded by %s → side-branch.\n",
			counterName(ev.absorbed), counterName(ev.by))
	}
	fmt.Fprintf(&b, "  → to flip the verdict, clear %s.\n", counterName(root))
	return b.String()
}
