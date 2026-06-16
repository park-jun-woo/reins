//ff:func feature=llm type=adapter control=sequence
//ff:what assertGrokCLISessionFallback — 한 세션 폴백 케이스를 검증하는 헬퍼: REINS_GROK_BIN을 비우고 REINS_GROK_SESSION을 설정/해제한 뒤 newGrokCLI("default")의 Session.Kind가 Stateless이고 Model이 빈값(default sentinel)인지 단언(TestNewGrokCLISessionFallback의 루프 본문 추출, Q4 회피).

package llm

import (
	"os"
	"testing"
)

// assertGrokCLISessionFallback applies one case: it blanks REINS_GROK_BIN, sets or
// unsets REINS_GROK_SESSION, and checks that newGrokCLI falls back to Stateless with
// an empty Model (the default sentinel). A misspelled value can never forge a session.
func assertGrokCLISessionFallback(t *testing.T, tc grokCLISessionFallbackCase) {
	t.Helper()
	t.Setenv("REINS_GROK_BIN", "")
	if tc.set {
		t.Setenv("REINS_GROK_SESSION", tc.val)
	} else {
		os.Unsetenv("REINS_GROK_SESSION")
	}
	g := newGrokCLI("default")
	if g.Session.Kind != Stateless {
		t.Fatalf("Session.Kind = %v, want Stateless", g.Session.Kind)
	}
	if g.Model != "" {
		t.Fatalf("Model = %q, want \"\" (default sentinel)", g.Model)
	}
}
