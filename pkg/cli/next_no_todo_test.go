//ff:func feature=cli type=helper control=sequence
//ff:what next가 아이템이 없을 때 남은 TODO가 없다고 알리는지 검증한다.

package cli

import (
	"path/filepath"
	"strings"
	"testing"
)

// TestNextNoTODO: with no items, next reports nothing remaining.
func TestNextNoTODO(t *testing.T) {
	dir := t.TempDir()
	session := filepath.Join(dir, "session.json")
	out := filepath.Join(dir, "out.jsonl")

	got := runCmd(t, stubDef{}, session, out, "", "next")
	if !strings.Contains(got, "no TODO items remaining") {
		t.Fatalf("next = %q", got)
	}
}
