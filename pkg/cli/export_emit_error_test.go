//ff:func feature=cli type=helper control=sequence level=error
//ff:what export가 sink의 Emit 실패 시 quest.Export 에러를 표면화하는지 검증한다(out 경로가 디렉터리라 OpenFile 실패).

package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/reins/pkg/quest"
)

// TestExportEmitError: export surfaces a quest.Export error when the sink's Emit
// fails. A terminal, not-yet-emitted item is seeded directly into the session, and
// the out path is an existing directory so newJSONLSink succeeds but Emit's OpenFile
// fails.
func TestExportEmitError(t *testing.T) {
	dir := t.TempDir()
	session := filepath.Join(dir, "session.json")
	s := &quest.Session{Version: 1, Items: []*quest.Item{{Key: "a", State: quest.PASS}}}
	if err := s.Save(session); err != nil {
		t.Fatalf("save: %v", err)
	}
	outDir := filepath.Join(dir, "outdir")
	if err := os.Mkdir(outDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	runCmdErr(t, stubDef{}, session, outDir, "", "export")
}
