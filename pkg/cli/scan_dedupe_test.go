//ff:func feature=cli type=helper control=sequence
//ff:what TestScanDedupe — 같은 입력으로 scan을 2회 실행하면 2회째는 전건 skip되어 총 아이템 수가 불변인지(래칫 무결성), 부분 중복 입력은 신규분만 추가하는지 검증한다.

package cli

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/park-jun-woo/reins/pkg/quest"
)

// TestScanDedupe: scanning the same input twice leaves the total item count
// unchanged (second run skips every duplicate), and a partially overlapping input
// adds only the new keys.
func TestScanDedupe(t *testing.T) {
	dir := t.TempDir()
	session := filepath.Join(dir, "session.json")
	out := filepath.Join(dir, "out.jsonl")

	got := runCmd(t, stubDef{}, session, out, "", "scan", "a", "b")
	if !strings.Contains(got, "seeded 2 item(s); 2 total") {
		t.Fatalf("first scan = %q", got)
	}

	// Identical re-scan: nothing added, total unchanged.
	got = runCmd(t, stubDef{}, session, out, "", "scan", "a", "b")
	if !strings.Contains(got, "seeded 0 item(s) (skipped 2 duplicate(s)); 2 total") {
		t.Fatalf("re-scan = %q", got)
	}
	s, err := quest.Load(session)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(s.Items) != 2 {
		t.Fatalf("total items = %d, want 2 (re-scan must not duplicate)", len(s.Items))
	}

	// Overlapping scan: only the new key lands.
	got = runCmd(t, stubDef{}, session, out, "", "scan", "b", "c")
	if !strings.Contains(got, "seeded 1 item(s) (skipped 1 duplicate(s)); 3 total") {
		t.Fatalf("overlap scan = %q", got)
	}
}
