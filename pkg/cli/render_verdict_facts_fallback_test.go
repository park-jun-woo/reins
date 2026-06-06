//ff:func feature=cli type=helper control=sequence
//ff:what TestRenderVerdictFactsFallback — Feedback이 없으면 각 Fact(규칙·위치·기대·실제)가 출력되는지 검증.

package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/park-jun-woo/reins/pkg/quest"
)

// TestRenderVerdictFactsFallback: without Feedback, each Fact is printed.
func TestRenderVerdictFactsFallback(t *testing.T) {
	it := &quest.Item{Key: "k", State: quest.TODO}
	v := quest.Verdict{
		Outcome: quest.OutFail,
		Facts: []quest.Fact{
			{Rule: "not-bad", Where: "body", Expected: "good", Actual: "bad"},
		},
	}
	var buf bytes.Buffer
	renderVerdict(&buf, "k", it, v)
	got := buf.String()
	if !strings.Contains(got, "k -> FAIL") {
		t.Fatalf("missing verdict line: %q", got)
	}
	if !strings.Contains(got, "not-bad") || !strings.Contains(got, "body") ||
		!strings.Contains(got, `expected="good"`) || !strings.Contains(got, `actual="bad"`) {
		t.Fatalf("fact not rendered: %q", got)
	}
}
