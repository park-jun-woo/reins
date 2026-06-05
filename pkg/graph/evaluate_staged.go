//ff:func feature=graph type=helper control=selection
//ff:what EvaluateStaged — §⑤ 단계적 평가로 G5(싼 검사 실패 시 fetch 생략)를 해소한다. ① tier0 서브그래프(Needs 없는 카운터+워런트)를 ground resolve 없이 평가 — 잔존 Fail이면 즉시 그 FAIL Verdict 반환(provider/resolver 절대 호출 안 됨, 네트워크 0회). ② tier0 clean이면 tier1이 선언한 ground를 provider+snapshot으로 1회씩 lazy resolve(캐시)해 ctx.Grounds에 주입 — resolve 에러는 terminal FAIL. ③ ground 주입된 ctx로 전체 그래프를 평가(기존 Evaluate: tier0/1 합산 verdict + §③ 공략집 렌더). snapshot은 호출측이 요청당 1개 생성(주입형 resolver로 결정론·테스트 네트워크-free).

package graph

import (
	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/ground"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// EvaluateStaged runs the §⑤ staged evaluation that resolves G5 (skip fetch when a
// cheap check fails). (1) It evaluates the tier-0 sub-graph (counters with no Needs,
// plus warrants) without resolving any ground; if a residual Fail remains it returns
// that FAIL Verdict immediately — the provider/resolver is never called and the
// network is touched zero times. (2) If tier-0 is clean, it resolves each ground
// declared by tier-1 counters once via provider+snapshot (lazy, cached) and injects
// them into ctx.Grounds — a resolve error is a terminal FAIL. (3) With the
// ground-populated ctx it evaluates the full graph (the existing Evaluate: combined
// tier-0/1 verdict plus the §③ walkthrough). The caller creates one snapshot per
// request (its injected resolver keeps evaluation deterministic and network-free in
// tests).
func (g *Graph) EvaluateStaged(ctx gate.Context, snap *ground.Snapshot, provider GroundProvider) quest.Verdict {
	// Tier 0: no-ground counters. Evaluated without resolving any ground.
	tier0 := g.subGraph("tier0", func(n *Node) bool {
		return len(n.needs) == 0
	})
	v0 := tier0.Evaluate(ctx)
	if v0.Outcome == quest.OutFail {
		// Residual Fail in tier 0 → short-circuit. Network untouched (G5).
		return v0
	}

	// Tier 0 clean → resolve the grounds tier-1 needs (lazy snapshot), then evaluate
	// the full graph.
	grounds, failV := g.resolveGrounds(ctx, snap, provider)
	if failV != nil {
		return *failV
	}
	ctx.Grounds = grounds
	return g.Evaluate(ctx)
}
