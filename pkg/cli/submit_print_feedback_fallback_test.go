//ff:func feature=cli type=helper control=sequence
//ff:what printSubmit이 Feedback이 비면 기존 Fact 폴백(규칙·위치·기대·실제 한 줄)으로 무회귀 출력하는지 검증한다.

package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/park-jun-woo/reins/pkg/quest"
)

// TestPrintSubmitFeedbackFallback: empty Feedback falls back to the Fact loop (no
// regression).
func TestPrintSubmitFeedbackFallback(t *testing.T) {
	var buf bytes.Buffer
	it := &quest.Item{Key: "d", State: quest.TODO}
	v := quest.Verdict{
		Outcome: quest.OutFail,
		Facts:   []quest.Fact{{Rule: "r", Where: "body", Expected: "good", Actual: "bad"}},
	}
	printSubmit(&buf, "d", it, v)
	got := buf.String()
	if !strings.Contains(got, `  - r: body expected="good" actual="bad"`) {
		t.Fatalf("Fact fallback missing: %q", got)
	}
}
