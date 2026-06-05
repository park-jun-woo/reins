//ff:func feature=graph type=helper control=sequence
//ff:what adaptRule 어댑터 증명 — 감싼 toulmin func가 gctx 미주입 시 (false,nil), gctx 타입 불일치 시 (false,nil), Check no-fire 시 (false,nil), Check fire 시 (true, fact)를 내고 fact.Rule이 ruleID로 스탬프되는지 검증한다.

package graph

import (
	"testing"

	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestAdaptRule(t *testing.T) {
	r := fireRule("dup", gate.LevelFail)
	fn := adaptRule(r, "stamped-id")

	t.Run("no gctx -> false", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ok, ev := fn(ctx, nil)
		if ok || ev != nil {
			t.Fatalf("no gctx: got (%v,%v) want (false,nil)", ok, ev)
		}
	})

	t.Run("wrong gctx type -> false", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set(gctxKey, "not-a-context")
		ok, ev := fn(ctx, nil)
		if ok || ev != nil {
			t.Fatalf("wrong type: got (%v,%v) want (false,nil)", ok, ev)
		}
	})

	t.Run("check no-fire -> false", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set(gctxKey, gate.Context{Submission: map[string]bool{}})
		ok, ev := fn(ctx, nil)
		if ok || ev != nil {
			t.Fatalf("no-fire: got (%v,%v) want (false,nil)", ok, ev)
		}
	})

	t.Run("check fire -> true, stamped fact", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set(gctxKey, gate.Context{Submission: map[string]bool{"dup": true}})
		ok, ev := fn(ctx, nil)
		if !ok {
			t.Fatalf("fire: got ok=false want true")
		}
		f, isFact := ev.(quest.Fact)
		if !isFact {
			t.Fatalf("fire: evidence type %T want quest.Fact", ev)
		}
		if f.Rule != "stamped-id" {
			t.Fatalf("fire: fact.Rule=%q want stamped-id", f.Rule)
		}
		if f.Where != "dup" {
			t.Fatalf("fire: fact.Where=%q want dup", f.Where)
		}
	})
}
