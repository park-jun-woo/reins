//ff:func feature=cli type=helper control=sequence level=error
//ff:what TestEvaluateAndApplyShortVerdict — Prepare 단락 verdict(SKIP)가 규칙 카탈로그를 거치지 않고 그대로 쓰이는지 검증.

package cli

import (
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/reins/pkg/quest"
)

// TestEvaluateAndApplyShortVerdict: a Prepare short-circuit verdict (SKIP) is used
// verbatim without invoking the rule catalog.
func TestEvaluateAndApplyShortVerdict(t *testing.T) {
	dir := t.TempDir()
	session := filepath.Join(dir, "session.json")
	out := filepath.Join(dir, "out.jsonl")

	s := quest.New()
	it := &quest.Item{Key: "a", State: quest.TODO}
	s.Items = append(s.Items, it)

	v, err := evaluateAndApply(stubDef{}, s, it, []byte("skip"), out, session)
	if err != nil {
		t.Fatalf("evaluateAndApply error: %v", err)
	}
	if v.Outcome != quest.OutSkip {
		t.Fatalf("outcome = %v, want SKIP", v.Outcome)
	}
}
