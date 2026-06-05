//ff:func feature=cli type=helper control=sequence
//ff:what status가 빈(부재) 세션에서 전부 0인 집계를 알리는지 검증한다.

package cli

import (
	"path/filepath"
	"strings"
	"testing"
)

// TestStatusEmpty: status on a fresh (absent) session reports an all-zero tally.
func TestStatusEmpty(t *testing.T) {
	dir := t.TempDir()
	session := filepath.Join(dir, "session.json")
	out := filepath.Join(dir, "out.jsonl")
	got := runCmd(t, stubDef{}, session, out, "", "status")
	if !strings.Contains(got, "TOTAL    0") {
		t.Fatalf("status = %q, want TOTAL 0", got)
	}
}
