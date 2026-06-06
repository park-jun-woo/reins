//ff:func feature=cli type=helper control=sequence
//ff:what TestRenderVerdictFeedbackPreferred — Verdict.Feedback가 있으면 (들여쓰기) 출력되고 Fact 루프를 건너뛰는지 검증.

package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/park-jun-woo/reins/pkg/quest"
)

// TestRenderVerdictFeedbackPreferred: when Verdict.Feedback is set, it is printed
// (indented) in place of the Fact loop.
func TestRenderVerdictFeedbackPreferred(t *testing.T) {
	it := &quest.Item{Key: "k", State: quest.TODO}
	v := quest.Verdict{
		Outcome:  quest.OutFail,
		Feedback: "line1\nline2\n",
		Facts:    []quest.Fact{{Rule: "r", Where: "w"}}, // must be ignored
	}
	var buf bytes.Buffer
	renderVerdict(&buf, "k", it, v)
	got := buf.String()
	if !strings.Contains(got, "k -> FAIL (state TODO)") {
		t.Fatalf("missing verdict line: %q", got)
	}
	if !strings.Contains(got, "  line1\n") || !strings.Contains(got, "  line2\n") {
		t.Fatalf("feedback not rendered: %q", got)
	}
	if strings.Contains(got, "r:") || strings.Contains(got, "expected") {
		t.Fatalf("fact loop should be skipped when Feedback set: %q", got)
	}
}
