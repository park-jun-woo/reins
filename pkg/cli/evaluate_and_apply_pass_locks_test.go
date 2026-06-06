//ff:func feature=cli type=helper control=sequence level=error
//ff:what TestEvaluateAndApplyPassLocks — flat Rules() 경로에서 통과 제출이 아이템을 PASS로 잠그고 세션·export 파일을 영속하는지 검증.

package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/reins/pkg/quest"
)

// TestEvaluateAndApplyPassLocks: a passing submission on the flat Rules() path locks
// the item PASS and persists the session and the export file.
func TestEvaluateAndApplyPassLocks(t *testing.T) {
	dir := t.TempDir()
	session := filepath.Join(dir, "session.json")
	out := filepath.Join(dir, "out.jsonl")

	s := quest.New()
	it := &quest.Item{Key: "a", State: quest.TODO}
	s.Items = append(s.Items, it)

	v, err := evaluateAndApply(stubDef{}, s, it, []byte("good"), out, session)
	if err != nil {
		t.Fatalf("evaluateAndApply error: %v", err)
	}
	if v.Outcome != quest.OutPass {
		t.Fatalf("outcome = %v, want PASS", v.Outcome)
	}
	if it.State != quest.PASS {
		t.Fatalf("state = %v, want PASS", it.State)
	}
	if _, err := os.Stat(session); err != nil {
		t.Fatalf("session not saved: %v", err)
	}
	data, err := os.ReadFile(out)
	if err != nil || len(data) == 0 {
		t.Fatalf("export not written: %v / %d bytes", err, len(data))
	}
}
