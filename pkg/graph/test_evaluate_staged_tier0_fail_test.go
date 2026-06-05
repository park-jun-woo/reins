//ff:func feature=graph type=helper control=sequence
//ff:what G5 수용 ① — tier-0 Fail 입력이 FAIL로 단락되고 resolver는 한 번도 호출되지 않음(네트워크 0)을 fake resolver 호출 카운터로 증명한다.

package graph

import (
	"testing"

	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/ground"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// TestEvaluateStagedTier0FailSkipsFetch — case ①: a tier-0 Fail short-circuits with
// FAIL and the resolver is never called (zero network).
func TestEvaluateStagedTier0FailSkipsFetch(t *testing.T) {
	g := stagedGraph()
	fr := &stagedFakeResolver{body: "ground-body"}
	snap := ground.NewSnapshot(fr)

	ctx := gate.Context{Submission: map[string]bool{"email-format": true}}
	v := g.EvaluateStaged(ctx, snap, stagedProvider)

	if v.Outcome != quest.OutFail {
		t.Fatalf("tier0 fail: outcome=%s want FAIL", v.Outcome)
	}
	if fr.fetchCalls != 0 {
		t.Fatalf("tier0 fail: resolver Fetch called %d times, want 0 (G5 violated)", fr.fetchCalls)
	}
}
