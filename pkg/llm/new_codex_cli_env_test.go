//ff:func feature=llm type=adapter control=sequence
//ff:what TestNewCodexCLIEnv — newCodexCLI가 REINS_CODEX_* env를 읽는지 검증: REINS_CODEX_BIN이 기본 "codex"를 덮어쓰고, REINS_CODEX_SESSION=="continue"가 Continue로 옵트인하며, model 토큰("default" 외)은 그대로 운반됨. 무서브프로세스.

package llm

import (
	"testing"
)

// TestNewCodexCLIEnv: newCodexCLI reads REINS_CODEX_* env. REINS_CODEX_BIN overrides
// the default "codex" binary, and REINS_CODEX_SESSION=="continue" opts into Continue;
// the model token is carried verbatim (only "default" is the empty-model sentinel).
func TestNewCodexCLIEnv(t *testing.T) {
	t.Setenv("REINS_CODEX_BIN", "/opt/codex")
	t.Setenv("REINS_CODEX_SESSION", "continue")

	c := newCodexCLI("gpt-5")
	if c.Bin != "/opt/codex" {
		t.Fatalf("Bin = %q, want /opt/codex (from REINS_CODEX_BIN)", c.Bin)
	}
	if c.Model != "gpt-5" {
		t.Fatalf("Model = %q, want gpt-5", c.Model)
	}
	if c.Session.Kind != Continue {
		t.Fatalf("Session.Kind = %v, want Continue (REINS_CODEX_SESSION=continue)", c.Session.Kind)
	}
}
