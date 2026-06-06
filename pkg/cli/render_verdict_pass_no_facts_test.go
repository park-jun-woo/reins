//ff:func feature=cli type=helper control=sequence
//ff:what TestRenderVerdictPassNoFacts — Fact 없는 PASS가 verdict 한 줄만 출력하는지 검증.

package cli

import (
	"bytes"
	"testing"

	"github.com/park-jun-woo/reins/pkg/quest"
)

// TestRenderVerdictPassNoFacts: a PASS with no facts prints only the verdict line.
func TestRenderVerdictPassNoFacts(t *testing.T) {
	it := &quest.Item{Key: "k", State: quest.PASS}
	var buf bytes.Buffer
	renderVerdict(&buf, "k", it, quest.Verdict{Outcome: quest.OutPass})
	got := buf.String()
	if got != "k -> PASS (state PASS)\n" {
		t.Fatalf("PASS render = %q", got)
	}
}
