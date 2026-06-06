//ff:func feature=graph type=helper control=iteration dimension=1
//ff:what Evaluate — §② 판독 모델. toulmin Evaluate(Trace:true)로 PASS 워런트를 평가하고 trace에서 활성(Activated) 카운터를 추출→각 노드의 reins Level과 Fact를 조인한다. Supersedes 적용: 활성 상류 카운터에 흡수된 하류 카운터를 집계에서 제외(잔존만). 잔존을 Level 집계 — 잔존 Fail≥1→FAIL / 잔존0→PASS / 잔존 전부 Review→REVIEW. Facts는 잔존(곁가지 제외)만. FAIL/REVIEW면 평가 중(이중평가 없음) §③ 공략집을 렌더해 Verdict.Feedback에 채운다(결정타+superseded 곁가지). 엔진 오류(사이클·panic)는 결정적 FAIL Fact로 환원. gate.verdict_from_levels 의미와 일관.

package graph

import (
	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

// Evaluate reads the verdict per the §② model: it evaluates the PASS warrant with
// toulmin (Trace on), extracts the activated counters from the trace, joins each
// node's reins Level and Fact, applies reins-side Supersedes (an active upstream
// counter excludes its absorbed downstream counters), then aggregates the
// remaining counters by Level: any remaining Fail ⇒ FAIL, none remaining ⇒ PASS,
// all remaining Review ⇒ REVIEW. Facts come only from the remaining (non
// side-branch) counters. An engine error (cycle/panic) reduces to a deterministic
// FAIL Fact. This is consistent with gate.verdictFromLevels.
func (g *Graph) Evaluate(gctx gate.Context) quest.Verdict {
	tctx := toulmin.NewContext()
	tctx.Set(gctxKey, gctx)

	results, err := g.tg.Evaluate(tctx, toulmin.EvalOption{Trace: true})
	if err != nil {
		return quest.Verdict{
			Outcome: quest.OutFail,
			Facts: []quest.Fact{{
				Rule:   "graph",
				Where:  "engine",
				Actual: err.Error(),
			}},
		}
	}

	// Collect activated counter nodes (with their Fact) across all warrant traces.
	activeCounters := g.collectActiveCounters(results)

	// Apply reins-side supersession: drop any active counter absorbed by an active
	// upstream counter.
	remaining := g.applySupersession(activeCounters)

	// Aggregate the remaining counters by Level.
	var facts []quest.Fact
	anyFail, anyReview := false, false
	for _, ac := range remaining {
		facts = append(facts, ac.fact)
		switch ac.node.level {
		case gate.LevelFail:
			anyFail = true
		case gate.LevelReview:
			anyReview = true
		}
	}
	v := verdictFromCounters(anyFail, anyReview, facts)

	// Render the agent-facing walkthrough in-place (no second evaluation): root
	// cause + superseded side-branches. Only meaningful for FAIL/REVIEW (PASS leaves
	// remaining empty ⇒ "").
	if v.Outcome != quest.OutPass {
		events := g.collectSupersessionEvents(activeCounters)
		v.Feedback = g.renderWalkthrough(string(v.Outcome), remaining, events)
		// RootCause = the same decisive counter the walkthrough selects (its
		// counterName = Fact.Rule). Structured exposure of an existing deterministic
		// choice; no new judgment logic.
		root, _ := g.selectRootCause(remaining)
		v.RootCause = counterName(root)
	}
	return v
}
