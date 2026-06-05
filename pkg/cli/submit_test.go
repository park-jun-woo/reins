//ff:func feature=cli type=helper control=sequence
//ff:what 통과 제출이 아이템을 PASS로 잠그는지 검증한다.

package cli

import (
	"path/filepath"
	"strings"
	"testing"
)

// TestSubmitPassLocks: a passing submission locks the item PASS.
func TestSubmitPassLocks(t *testing.T) {
	dir := t.TempDir()
	session := filepath.Join(dir, "session.json")
	out := filepath.Join(dir, "out.jsonl")
	runCmd(t, stubDef{}, session, out, "", "scan", "a")

	got := runCmd(t, stubDef{}, session, out, "good", "submit", "--key", "a")
	if !strings.Contains(got, "a -> PASS") || !strings.Contains(got, "state PASS") {
		t.Fatalf("submit = %q", got)
	}
}
