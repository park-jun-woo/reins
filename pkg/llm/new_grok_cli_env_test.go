//ff:func feature=llm type=adapter control=sequence
//ff:what TestNewGrokCLIEnv — newGrokCLI가 REINS_GROK_* env를 읽는지 검증: REINS_GROK_BIN이 기본 "grok"을 덮어쓰고, REINS_GROK_SESSION=="continue"가 Continue로 옵트인하며, model 토큰("default" 외)은 그대로 운반됨. 무서브프로세스.

package llm

import (
	"testing"
)

// TestNewGrokCLIEnv: newGrokCLI reads REINS_GROK_* env. REINS_GROK_BIN overrides
// the default "grok" binary, and REINS_GROK_SESSION=="continue" opts into Continue;
// the model token is carried verbatim (only "default" is the empty-model sentinel).
func TestNewGrokCLIEnv(t *testing.T) {
	t.Setenv("REINS_GROK_BIN", "/opt/grok")
	t.Setenv("REINS_GROK_SESSION", "continue")

	g := newGrokCLI("grok-4")
	if g.Bin != "/opt/grok" {
		t.Fatalf("Bin = %q, want /opt/grok (from REINS_GROK_BIN)", g.Bin)
	}
	if g.Model != "grok-4" {
		t.Fatalf("Model = %q, want grok-4", g.Model)
	}
	if g.Session.Kind != Continue {
		t.Fatalf("Session.Kind = %v, want Continue (REINS_GROK_SESSION=continue)", g.Session.Kind)
	}
}
