//ff:func feature=graph type=helper control=iteration dimension=1
//ff:what 결정론·멱등 증명 — 같은 그래프를 같은 gate.Context로 두 번 Evaluate하면 동일한 Verdict(Outcome + Facts)가 나오는지 검증한다. 카운터 발동·supersession이 섞인 케이스를 써서 trace 판독·집계가 ctx에만 의존하고 평가 간 상태가 새지 않음을(고정 qualifier·순수 ctx → 결정론) 보장한다.

package graph

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/reins/pkg/gate"
)

func TestGraphEvaluateIdempotent(t *testing.T) {
	rFmt := fireRule("fmt", gate.LevelFail)
	rHolder := fireRule("holder", gate.LevelFail)
	rFree := fireRule("free", gate.LevelReview)

	g := NewGraph("idempotent")
	pass := g.Warrant(alwaysTrue)
	fmtN := g.Counter(rFmt, gate.LevelFail).Attacks(pass)
	holder := g.Counter(rHolder, gate.LevelFail).Attacks(pass)
	free := g.Counter(rFree, gate.LevelReview).Attacks(pass)
	fmtN.Supersedes(holder, free)
	free.Supersedes(holder)

	ctx := gate.Context{Submission: map[string]bool{"holder": true, "free": true}}

	first := g.Evaluate(ctx)
	for i := 0; i < 3; i++ {
		again := g.Evaluate(ctx)
		if !reflect.DeepEqual(first, again) {
			t.Fatalf("evaluate not idempotent: first=%+v again=%+v", first, again)
		}
	}
}
