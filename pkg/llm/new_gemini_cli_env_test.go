//ff:func feature=llm type=adapter control=sequence
//ff:what TestNewGeminiCLIEnv — newGeminiCLI가 REINS_GEMINI_* env를 읽는지 검증: REINS_GEMINI_BIN이 기본 "gemini"를 덮어쓰고, REINS_GEMINI_SESSION=="continue"가 Continue로 옵트인하며, model 토큰("default" 외)은 그대로 운반됨. 무서브프로세스.

package llm

import (
	"testing"
)

// TestNewGeminiCLIEnv: newGeminiCLI reads REINS_GEMINI_* env. REINS_GEMINI_BIN
// overrides the default "gemini" binary, and REINS_GEMINI_SESSION=="continue" opts
// into Continue; the model token is carried verbatim (only "default" is the
// empty-model sentinel).
func TestNewGeminiCLIEnv(t *testing.T) {
	t.Setenv("REINS_GEMINI_BIN", "/opt/gemini")
	t.Setenv("REINS_GEMINI_SESSION", "continue")

	c := newGeminiCLI("gemini-2.5-pro")
	if c.Bin != "/opt/gemini" {
		t.Fatalf("Bin = %q, want /opt/gemini (from REINS_GEMINI_BIN)", c.Bin)
	}
	if c.Model != "gemini-2.5-pro" {
		t.Fatalf("Model = %q, want gemini-2.5-pro", c.Model)
	}
	if c.Session.Kind != Continue {
		t.Fatalf("Session.Kind = %v, want Continue (REINS_GEMINI_SESSION=continue)", c.Session.Kind)
	}
}
