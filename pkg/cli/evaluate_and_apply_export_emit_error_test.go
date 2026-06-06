//ff:func feature=cli type=helper control=sequence level=error
//ff:what TestEvaluateAndApplyExportEmitError — newJSONLSink는 성공(부모 디렉터리 존재)하지만 out 경로가 디렉터리라 Emit이 실패하는 Export 에러 분기를 자극한다(sink 생성 분기와 구분).

package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/reins/pkg/quest"
)

// TestEvaluateAndApplyExportEmitError: newJSONLSink succeeds (parent dir exists) but
// Emit fails because the out path is itself a directory — exercising the Export
// error branch (distinct from the sink-construction branch).
func TestEvaluateAndApplyExportEmitError(t *testing.T) {
	dir := t.TempDir()
	session := filepath.Join(dir, "session.json")
	// out is an existing directory: MkdirAll(parent) succeeds, OpenFile(out) fails.
	out := filepath.Join(dir, "outdir")
	if err := os.MkdirAll(out, 0o755); err != nil {
		t.Fatal(err)
	}

	s := quest.New()
	it := &quest.Item{Key: "a", State: quest.TODO}
	s.Items = append(s.Items, it)

	if _, err := evaluateAndApply(stubDef{}, s, it, []byte("good"), out, session); err == nil {
		t.Fatal("evaluateAndApply = nil error, want Export emit error")
	}
}
