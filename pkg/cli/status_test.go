//ff:func feature=cli type=helper control=sequence
//ff:what status가 PASS 후 상태별 집계와 TOTAL을 알리는지 검증한다.

package cli

import (
	"path/filepath"
	"strings"
	"testing"
)

// TestStatusTally: status reports the per-state tally and TOTAL after a PASS.
func TestStatusTally(t *testing.T) {
	dir := t.TempDir()
	session := filepath.Join(dir, "session.json")
	out := filepath.Join(dir, "out.jsonl")
	runCmd(t, stubDef{}, session, out, "", "scan", "a", "b")
	runCmd(t, stubDef{}, session, out, "good", "submit", "--key", "a")

	got := runCmd(t, stubDef{}, session, out, "", "status")
	if !strings.Contains(got, "PASS     1") || !strings.Contains(got, "TODO     1") {
		t.Fatalf("status tally = %q", got)
	}
	if !strings.Contains(got, "TOTAL    2") {
		t.Fatalf("status total = %q", got)
	}
}
