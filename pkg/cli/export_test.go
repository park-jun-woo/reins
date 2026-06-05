//ff:func feature=cli type=helper control=sequence
//ff:what export가 종단 아이템을 JSONL로 쓰고 멱등인지(2회차 0건) 검증한다.

package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestExportTerminal: export writes terminal items to JSONL and is idempotent — a
// second export emits 0 new records.
func TestExportTerminal(t *testing.T) {
	dir := t.TempDir()
	session := filepath.Join(dir, "session.json")
	out := filepath.Join(dir, "out.jsonl")
	runCmd(t, stubDef{}, session, out, "", "scan", "a")
	runCmd(t, stubDef{}, session, out, "good", "submit", "--key", "a")

	// submit already swept the PASS item, so a follow-up export adds nothing new.
	got := runCmd(t, stubDef{}, session, out, "", "export")
	if !strings.Contains(got, "exported 0 new record(s)") {
		t.Fatalf("export = %q", got)
	}
	b, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("read out: %v", err)
	}
	if !strings.Contains(string(b), `"key":"a"`) {
		t.Fatalf("out file = %q", b)
	}
}
