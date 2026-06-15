//ff:func feature=cli type=helper control=sequence level=error
//ff:what TestApplyVerdictPassLocks — PASS verdict 적용 시 아이템이 PASS로 잠기고 세션·export(JSONL)가 영속되는지 단언. applyVerdict는 게이트가 건넨 PASS를 그대로 적용·emit하는 단일 지점임을 고정한다.

package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/reins/pkg/quest"
)

// TestApplyVerdictPassLocks: applying a PASS verdict locks the item PASS and persists
// both the session and the export (JSONL) file.
func TestApplyVerdictPassLocks(t *testing.T) {
	dir := t.TempDir()
	session := filepath.Join(dir, "session.json")
	out := filepath.Join(dir, "out.jsonl")

	s := quest.New()
	it := &quest.Item{Key: "a", State: quest.TODO}
	s.Items = append(s.Items, it)

	v := quest.Verdict{Outcome: quest.OutPass}
	if err := applyVerdict(s, it, v, out, session); err != nil {
		t.Fatalf("applyVerdict error: %v", err)
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
