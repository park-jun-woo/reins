//ff:func feature=cli type=helper control=sequence level=error
//ff:what TestApplyVerdictFailRatchets — FAIL verdict 적용 시 it.Tries가 1 증가하고(잠금 없음, State는 TODO 유지) 세션이 디스크에 Save되어 재로드 가능한지 단언. verdict가 래칫을 바꾸는 단일 지점의 FAIL 분기를 고정한다.

package cli

import (
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/reins/pkg/quest"
)

// TestApplyVerdictFailRatchets: applying a FAIL verdict bumps it.Tries by one (no
// lock — the item stays TODO) and persists the session so it reloads from disk.
func TestApplyVerdictFailRatchets(t *testing.T) {
	dir := t.TempDir()
	session := filepath.Join(dir, "session.json")
	out := filepath.Join(dir, "out.jsonl")

	s := quest.New()
	it := &quest.Item{Key: "a", State: quest.TODO}
	s.Items = append(s.Items, it)

	v := quest.Verdict{Outcome: quest.OutFail, RootCause: "r1"}
	if err := applyVerdict(s, it, v, out, session); err != nil {
		t.Fatalf("applyVerdict error: %v", err)
	}
	if it.Tries != 1 {
		t.Fatalf("tries = %d, want 1", it.Tries)
	}
	if it.State != quest.TODO {
		t.Fatalf("state = %v, want TODO (FAIL must not lock before MaxTries)", it.State)
	}

	reloaded, err := quest.Load(session)
	if err != nil {
		t.Fatalf("session not reloadable: %v", err)
	}
	if len(reloaded.Items) != 1 || reloaded.Items[0].Tries != 1 {
		t.Fatalf("reloaded tries = %v, want 1", reloaded.Items)
	}
}
