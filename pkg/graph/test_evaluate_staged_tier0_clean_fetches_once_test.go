//ff:func feature=graph type=helper control=iteration dimension=1
//ff:what G5 수용 ② — clean tier-0 입력이 ground를 resolve하되 두 tier-1이 같은 URL을 읽어도 Fetch가 정확히 1회(스냅샷 캐시)이고, tier-1 verdict가 resolve된 본문을 반영(FAIL + Actual=본문)함을 증명한다. Facts를 훑어 resolve된 ground가 실린 Fact를 찾는다.

package graph

import (
	"testing"

	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/ground"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// TestEvaluateStagedTier0CleanFetchesOnce — case ②: a clean tier-0 resolves the
// ground (Fetch exactly once despite two tier-1 reads of the same URL) and the
// tier-1 verdict reflects the resolved body.
func TestEvaluateStagedTier0CleanFetchesOnce(t *testing.T) {
	g := stagedGraph()
	fr := &stagedFakeResolver{body: "ground-body"}
	snap := ground.NewSnapshot(fr)

	// email-format does NOT fire → tier 0 clean.
	ctx := gate.Context{Submission: map[string]bool{}}
	v := g.EvaluateStaged(ctx, snap, stagedProvider)

	if fr.fetchCalls != 1 {
		t.Fatalf("tier0 clean: resolver Fetch called %d times, want exactly 1 (snapshot cache)", fr.fetchCalls)
	}
	// Both tier-1 rules fire on the non-empty body → FAIL with the ground reflected.
	if v.Outcome != quest.OutFail {
		t.Fatalf("tier0 clean: outcome=%s want FAIL (tier-1 fired on ground)", v.Outcome)
	}
	found := false
	for _, f := range v.Facts {
		if f.Actual == "ground-body" {
			found = true
		}
	}
	if !found {
		t.Fatalf("tier0 clean: no fact carried the resolved ground; facts=%v", v.Facts)
	}
}
