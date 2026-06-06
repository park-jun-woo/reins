//ff:func feature=cli type=helper control=sequence level=error
//ff:what TestEvaluateAndApplyPrepareError — Prepare 에러가 그대로 반환되고(apply 없음) 아이템이 TODO로 남는지 검증.

package cli

import (
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/reins/pkg/quest"
)

// TestEvaluateAndApplyPrepareError: a Prepare error is returned directly (no apply).
func TestEvaluateAndApplyPrepareError(t *testing.T) {
	dir := t.TempDir()
	session := filepath.Join(dir, "session.json")
	out := filepath.Join(dir, "out.jsonl")

	s := quest.New()
	it := &quest.Item{Key: "a", State: quest.TODO}
	s.Items = append(s.Items, it)

	if _, err := evaluateAndApply(errDef{prepareErr: true}, s, it, []byte("x"), out, session); err == nil {
		t.Fatal("evaluateAndApply = nil error, want Prepare error")
	}
	if it.State != quest.TODO {
		t.Fatalf("state = %v, want TODO (no apply on Prepare error)", it.State)
	}
}
