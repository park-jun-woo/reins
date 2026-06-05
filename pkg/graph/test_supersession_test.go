//ff:func feature=graph type=helper control=iteration dimension=1
//ff:what supersession 판독 증명 — free(Review)가 holder(Fail)를 Supersedes할 때, 둘 다 발동해도 holder가 집계에서 제외되어 잔존=free만 → REVIEW가 나오는지(엣지0이면 FAIL일 것을 곁가지 흡수로 뒤집음) + 상류 fmt(Fail)가 활성이면 fmt는 잔존해 FAIL을 유지하는지 검증한다.

package graph

import (
	"testing"

	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
)

func TestSupersessionAbsorbsSideBranch(t *testing.T) {
	rFmt := fireRule("fmt", gate.LevelFail)
	rHolder := fireRule("holder", gate.LevelFail)
	rFree := fireRule("free", gate.LevelReview)

	build := func() *Graph {
		g := NewGraph("supersede")
		pass := g.Warrant(alwaysTrue)
		fmtN := g.Counter(rFmt, gate.LevelFail).Attacks(pass)
		holder := g.Counter(rHolder, gate.LevelFail).Attacks(pass)
		free := g.Counter(rFree, gate.LevelReview).Attacks(pass)
		fmtN.Supersedes(holder, free) // fmt absorbs holder + free
		free.Supersedes(holder)       // free absorbs holder
		return g
	}

	cases := []struct {
		name  string
		fired map[string]bool
		want  quest.Outcome
	}{
		{"free absorbs holder -> review", map[string]bool{"holder": true, "free": true}, quest.OutReview},
		{"fmt remains -> fail", map[string]bool{"fmt": true, "holder": true, "free": true}, quest.OutFail},
		{"holder alone -> fail", map[string]bool{"holder": true}, quest.OutFail},
		{"none -> pass", map[string]bool{}, quest.OutPass},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			g := build()
			v := g.Evaluate(gate.Context{Submission: c.fired})
			if v.Outcome != c.want {
				t.Fatalf("outcome=%s want=%s (facts=%v)", v.Outcome, c.want, factWheres(v))
			}
		})
	}
}
