//ff:func feature=llm type=adapter control=sequence level=error
//ff:what TestGeminiCompleteNoKey — GEMINI_API_KEY 미설정 시 네트워크 호출 전에 에러로 단락되는지 검증.

package llm

import (
	"testing"
)

// TestGeminiCompleteNoKey: a missing GEMINI_API_KEY errors before any network call.
func TestGeminiCompleteNoKey(t *testing.T) {
	t.Setenv("GEMINI_API_KEY", "")
	g := Gemini{Model: "gemini-1.5-pro"}
	got, err := g.Complete("sys", "usr")
	if err == nil {
		t.Fatalf("Complete = %q, want missing-key error", got)
	}
	if got != "" {
		t.Fatalf("result = %q, want empty on error", got)
	}
}
