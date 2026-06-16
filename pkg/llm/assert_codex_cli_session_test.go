//ff:func feature=llm type=adapter control=sequence
//ff:what assertCodexCLISessionFallback — 한 세션 폴백 케이스를 검증하는 헬퍼: REINS_CODEX_BIN을 비우고 REINS_CODEX_SESSION을 설정/해제한 뒤 newCodexCLI("default")의 Session.Kind가 Stateless이고 Model이 빈값(default sentinel)인지 단언(TestNewCodexCLISessionFallback의 루프 본문 추출, Q4 회피).

package llm

import (
	"os"
	"testing"
)

// assertCodexCLISessionFallback applies one case: it blanks REINS_CODEX_BIN, sets or
// unsets REINS_CODEX_SESSION, and checks that newCodexCLI falls back to Stateless with
// an empty Model (the default sentinel). A misspelled value can never forge a session.
func assertCodexCLISessionFallback(t *testing.T, tc codexCLISessionFallbackCase) {
	t.Helper()
	t.Setenv("REINS_CODEX_BIN", "")
	if tc.set {
		t.Setenv("REINS_CODEX_SESSION", tc.val)
	} else {
		os.Unsetenv("REINS_CODEX_SESSION")
	}
	c := newCodexCLI("default")
	if c.Session.Kind != Stateless {
		t.Fatalf("Session.Kind = %v, want Stateless", c.Session.Kind)
	}
	if c.Model != "" {
		t.Fatalf("Model = %q, want \"\" (default sentinel)", c.Model)
	}
}
