//ff:func feature=llm type=loader control=sequence level=error
//ff:what TestLoadAPIKeyMissing — env var가 미설정이면 에러를 반환하고 키는 빈 문자열인지 검증.

package llm

import (
	"testing"
)

// TestLoadAPIKeyMissing: an env var that is unset yields an error.
func TestLoadAPIKeyMissing(t *testing.T) {
	t.Setenv("GEMINI_API_KEY", "")
	got, err := loadAPIKey("gemini")
	if err == nil {
		t.Fatalf("loadAPIKey = %q, want error when env unset", got)
	}
	if got != "" {
		t.Fatalf("key = %q, want empty on error", got)
	}
}
