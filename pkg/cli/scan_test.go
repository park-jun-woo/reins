//ff:func feature=cli type=helper control=sequence
//ff:what scan이 arg 하나당 아이템을 시드하고 누적 합계를 알리는지 검증한다.

package cli

import (
	"path/filepath"
	"strings"
	"testing"
)

// TestScanSeeds: scan seeds one item per arg and reports the running total.
func TestScanSeeds(t *testing.T) {
	dir := t.TempDir()
	session := filepath.Join(dir, "session.json")
	out := filepath.Join(dir, "out.jsonl")

	got := runCmd(t, stubDef{}, session, out, "", "scan", "a", "b")
	if !strings.Contains(got, "seeded 2 item(s); 2 total") {
		t.Fatalf("scan = %q", got)
	}
}
