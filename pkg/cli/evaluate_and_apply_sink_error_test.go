//ff:func feature=cli type=helper control=sequence level=error
//ff:what TestEvaluateAndApplySinkError — 쓰기 불가 export 경로(부모가 파일)가 verdict 계산 후 sink 생성 에러를 표면화하는지 검증(OS별 메시지 차이에 관대).

package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/park-jun-woo/reins/pkg/quest"
)

// TestEvaluateAndApplySinkError: an unwritable export path surfaces an error after
// the verdict is computed.
func TestEvaluateAndApplySinkError(t *testing.T) {
	dir := t.TempDir()
	session := filepath.Join(dir, "session.json")
	// A path whose parent is a file (not a dir) makes the sink mkdir/open fail.
	badParent := filepath.Join(dir, "afile")
	if err := os.WriteFile(badParent, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(badParent, "nested", "out.jsonl")

	s := quest.New()
	it := &quest.Item{Key: "a", State: quest.TODO}
	s.Items = append(s.Items, it)

	if _, err := evaluateAndApply(stubDef{}, s, it, []byte("good"), out, session); err == nil {
		t.Fatal("evaluateAndApply = nil error, want sink error")
	} else if !strings.Contains(err.Error(), "afile") && !strings.Contains(strings.ToLower(err.Error()), "not a directory") {
		// Tolerant: just ensure an error happened; the message varies by OS.
		t.Logf("sink error (ok): %v", err)
	}
}
