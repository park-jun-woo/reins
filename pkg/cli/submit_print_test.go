//ff:func feature=cli type=helper control=sequence
//ff:what printSubmit이 Facts가 없으면 verdict 한 줄만 출력하는지 검증한다.

package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/park-jun-woo/reins/pkg/quest"
)

// TestPrintSubmitVerdictLine: with no Facts, printSubmit emits only the verdict line.
func TestPrintSubmitVerdictLine(t *testing.T) {
	var buf bytes.Buffer
	it := &quest.Item{Key: "a", State: quest.PASS}
	printSubmit(&buf, "a", it, quest.Verdict{Outcome: quest.OutPass})
	got := buf.String()
	if !strings.Contains(got, "a -> PASS (state PASS)") {
		t.Fatalf("printSubmit = %q", got)
	}
	if strings.Contains(got, "  - ") {
		t.Fatalf("no Facts expected, got %q", got)
	}
}
