//ff:func feature=cli type=helper control=sequence level=error
//ff:what TestEvaluateAndApplySessionSaveError — 쓰기 불가 세션 경로(부모가 파일)가 verdict 계산·적용 후 Save 에러를 표면화하는지 검증.

package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/reins/pkg/quest"
)

// TestEvaluateAndApplySessionSaveError: an unwritable session path surfaces the
// save error after the verdict is computed and applied.
func TestEvaluateAndApplySessionSaveError(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "out.jsonl")
	// A session path under a file (not a dir) makes Save fail.
	badParent := filepath.Join(dir, "afile")
	if err := os.WriteFile(badParent, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	session := filepath.Join(badParent, "nested", "session.json")

	s := quest.New()
	it := &quest.Item{Key: "a", State: quest.TODO}
	s.Items = append(s.Items, it)

	if _, err := evaluateAndApply(stubDef{}, s, it, []byte("good"), out, session); err == nil {
		t.Fatal("evaluateAndApply = nil error, want session save error")
	}
}
