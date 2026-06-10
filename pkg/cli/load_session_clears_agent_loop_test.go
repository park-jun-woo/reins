//ff:func feature=cli type=helper control=sequence level=error
//ff:what TestLoadSessionClearsAgentLoop — MetaAgentLoop 플래그가 박힌 세션 파일(kill된 agent 프로세스의 잔류)을 loadSession이 로드 직후 자가 치유(플래그 삭제)하는지, 다른 Meta 키는 보존하는지 검증한다.

package cli

import (
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/reins/pkg/quest"
)

// TestLoadSessionClearsAgentLoop: a session file persisted with the MetaAgentLoop
// flag set (residue of a killed agent process) is self-healed on load — the flag is
// gone, while other Meta keys survive.
func TestLoadSessionClearsAgentLoop(t *testing.T) {
	path := filepath.Join(t.TempDir(), "session.json")
	stale := quest.New()
	stale.SetMeta(quest.MetaAgentLoop, true)
	stale.SetMeta("keep", "me")
	if err := stale.Save(path); err != nil {
		t.Fatalf("save: %v", err)
	}

	s, err := loadSession(path)
	if err != nil {
		t.Fatalf("loadSession: %v", err)
	}
	if _, ok := s.GetMeta(quest.MetaAgentLoop); ok {
		t.Fatal("MetaAgentLoop survived load; want it cleared (self-heal)")
	}
	if v, ok := s.GetMeta("keep"); !ok || v != "me" {
		t.Fatalf("other Meta key lost: keep=%v ok=%v", v, ok)
	}
}
