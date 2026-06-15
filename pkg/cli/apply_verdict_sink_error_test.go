//ff:func feature=cli type=helper control=sequence level=error
//ff:what TestApplyVerdictSinkError — 세션 Save는 성공하되 export sink 생성이 실패하는 outPath(정규 파일 하위 경로)에서 applyVerdict가 newJSONLSink 에러를 전파하는지 단언. Save와 sink-open 두 실패 분기를 분리해 고정한다.

package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/reins/pkg/quest"
)

// TestApplyVerdictSinkError: the session saves fine, but the export sink cannot be
// created (outPath nested under a regular file), so applyVerdict propagates the
// newJSONLSink error — distinct from the Save-failure branch.
func TestApplyVerdictSinkError(t *testing.T) {
	dir := t.TempDir()
	session := filepath.Join(dir, "session.json")
	badParent := filepath.Join(dir, "afile")
	if err := os.WriteFile(badParent, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(badParent, "nested", "out.jsonl")

	s := quest.New()
	it := &quest.Item{Key: "a", State: quest.TODO}
	s.Items = append(s.Items, it)

	v := quest.Verdict{Outcome: quest.OutFail}
	if err := applyVerdict(s, it, v, out, session); err == nil {
		t.Fatal("err = nil, want newJSONLSink error to propagate")
	}
	if _, statErr := os.Stat(session); statErr != nil {
		t.Fatalf("session should have been saved before the sink error: %v", statErr)
	}
}
