//ff:func feature=cli type=helper control=sequence level=error
//ff:what TestApplyVerdictSaveError — 쓰기 불가 sessionPath(정규 파일 하위 경로)에서 applyVerdict가 Save 실패를 그대로 전파하는지 단언. 정합성 위협(영속화 실패)은 침묵 금지.

package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/reins/pkg/quest"
)

// TestApplyVerdictSaveError: when the session path is unwritable (nested under a
// regular file), applyVerdict propagates the Save error rather than swallowing it —
// a persistence failure must never be silent.
func TestApplyVerdictSaveError(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "out.jsonl")
	badParent := filepath.Join(dir, "afile")
	if err := os.WriteFile(badParent, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	session := filepath.Join(badParent, "nested", "session.json")

	s := quest.New()
	it := &quest.Item{Key: "a", State: quest.TODO}
	s.Items = append(s.Items, it)

	v := quest.Verdict{Outcome: quest.OutFail}
	if err := applyVerdict(s, it, v, out, session); err == nil {
		t.Fatal("err = nil, want Save error to propagate")
	}
}
