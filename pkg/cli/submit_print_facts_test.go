//ff:func feature=cli type=helper control=sequence
//ff:what printSubmit이 Facts를 들여쓴 자기교정 피드백 줄로 렌더하는지 검증한다.

package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/park-jun-woo/reins/pkg/quest"
)

// TestPrintSubmitFacts: Facts are rendered as indented self-correction feedback lines.
func TestPrintSubmitFacts(t *testing.T) {
	var buf bytes.Buffer
	it := &quest.Item{Key: "b", State: quest.TODO}
	v := quest.Verdict{
		Outcome: quest.OutFail,
		Facts:   []quest.Fact{{Rule: "not-bad", Where: "body", Expected: "good", Actual: "bad"}},
	}
	printSubmit(&buf, "b", it, v)
	got := buf.String()
	if !strings.Contains(got, "b -> FAIL (state TODO)") {
		t.Fatalf("verdict line = %q", got)
	}
	if !strings.Contains(got, `  - not-bad: body expected="good" actual="bad"`) {
		t.Fatalf("fact line = %q", got)
	}
}
