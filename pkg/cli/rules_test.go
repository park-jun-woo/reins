//ff:func feature=cli type=helper control=sequence
//ff:what rules가 게이트 카탈로그(레벨·ID·설명)를 출력하는지 검증한다.

package cli

import (
	"path/filepath"
	"strings"
	"testing"
)

// TestRulesCatalog: rules prints the gate's catalog (level, id, description).
func TestRulesCatalog(t *testing.T) {
	dir := t.TempDir()
	session := filepath.Join(dir, "session.json")
	out := filepath.Join(dir, "out.jsonl")

	got := runCmd(t, stubDef{}, session, out, "", "rules")
	if !strings.Contains(got, "not-bad") || !strings.Contains(got, "submission must not be bad") {
		t.Fatalf("rules = %q", got)
	}
}
