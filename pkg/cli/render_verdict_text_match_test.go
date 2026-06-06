//ff:func feature=cli type=helper control=sequence
//ff:what TestRenderVerdictText — renderVerdictText가 renderVerdict와 바이트 동일 출력을 string으로 반환하는지(사람 피드백=모델 피드백) 검증.

package cli

import (
	"bytes"
	"testing"

	"github.com/park-jun-woo/reins/pkg/quest"
)

// TestRenderVerdictText matches renderVerdict byte-for-byte (human feedback == model
// feedback).
func TestRenderVerdictText(t *testing.T) {
	it := &quest.Item{Key: "k", State: quest.TODO}
	v := quest.Verdict{Outcome: quest.OutFail, Feedback: "why\n"}
	var buf bytes.Buffer
	renderVerdict(&buf, "k", it, v)
	if got := renderVerdictText("k", it, v); got != buf.String() {
		t.Fatalf("renderVerdictText = %q, want %q", got, buf.String())
	}
}
