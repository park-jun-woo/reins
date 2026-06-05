//ff:func feature=cli type=helper control=sequence
//ff:what printSubmit이 Verdict.Feedback이 있으면 그 공략집을 들여써 출력하고 Facts 루프는 건너뛰는지(Feedback 우선) 검증한다.

package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/park-jun-woo/reins/pkg/quest"
)

// TestPrintSubmitFeedbackPreferred: when Feedback is set, it is printed (indented)
// and the per-Fact loop is skipped.
func TestPrintSubmitFeedbackPreferred(t *testing.T) {
	var buf bytes.Buffer
	it := &quest.Item{Key: "c", State: quest.TODO}
	v := quest.Verdict{
		Outcome:  quest.OutFail,
		Facts:    []quest.Fact{{Rule: "should-not-print", Where: "x"}},
		Feedback: "FAIL. root cause = email-format (remaining active FAIL, upstream).\n  → to flip the verdict, clear email-format.\n",
	}
	printSubmit(&buf, "c", it, v)
	got := buf.String()
	if !strings.Contains(got, "c -> FAIL (state TODO)") {
		t.Fatalf("verdict line = %q", got)
	}
	if !strings.Contains(got, "root cause = email-format") {
		t.Fatalf("Feedback not printed: %q", got)
	}
	if !strings.Contains(got, "to flip the verdict, clear email-format") {
		t.Fatalf("Feedback action not printed: %q", got)
	}
	if strings.Contains(got, "should-not-print") {
		t.Fatalf("Facts loop must be skipped when Feedback is set: %q", got)
	}
}
