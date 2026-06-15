//ff:func feature=cli type=helper control=sequence level=error
//ff:what TestApplyVerdictExportError — sink 생성은 성공하되 export Emit이 실패하는 outPath(디렉터리 경로)에서 applyVerdict가 exportAndSave 에러를 전파하는지 단언. PASS로 아이템을 terminal화해 Export가 실제 Emit을 시도하게 한 뒤 그 실패 분기를 고정한다.

package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/reins/pkg/quest"
)

// TestApplyVerdictExportError: the sink is created fine, but the export Emit fails
// because outPath is a directory (OpenFile on a dir errors). A PASS verdict makes the
// item terminal so Export actually tries to emit it, exercising the exportAndSave
// error branch — distinct from the Save and sink-open branches.
func TestApplyVerdictExportError(t *testing.T) {
	dir := t.TempDir()
	session := filepath.Join(dir, "session.json")
	out := filepath.Join(dir, "outdir")
	if err := os.Mkdir(out, 0o755); err != nil {
		t.Fatal(err)
	}

	s := quest.New()
	it := &quest.Item{Key: "a", State: quest.TODO}
	s.Items = append(s.Items, it)

	v := quest.Verdict{Outcome: quest.OutPass}
	if err := applyVerdict(s, it, v, out, session); err == nil {
		t.Fatal("err = nil, want exportAndSave (Emit) error to propagate")
	}
}
