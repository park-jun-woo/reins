//ff:func feature=llm type=adapter control=sequence level=error
//ff:what TestGeminiCompleteConnError — 도달 불가 주소면 요청 에러인지 검증(실제 네트워크 성공 없이 분기 커버).

package llm

import (
	"testing"
)

// TestGeminiCompleteConnError: an unreachable address yields a request error.
func TestGeminiCompleteConnError(t *testing.T) {
	t.Setenv("GEMINI_API_KEY", "gk")
	g := Gemini{Model: "m", BaseURL: "http://127.0.0.1:1"}
	if _, err := g.Complete("a", "b"); err == nil {
		t.Fatal("Complete = nil error, want request error")
	}
}
