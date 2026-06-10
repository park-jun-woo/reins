//ff:func feature=llm type=adapter control=sequence level=error
//ff:what TestGeminiCompleteKeyNotInError — 연결 실패의 *url.Error가 전체 URL을 포함하므로, 키가 헤더로만 전달되어 에러 문자열에 API 키 값이 절대 나타나지 않는지 검증한다(키 누출 차단 게이트).

package llm

import (
	"strings"
	"testing"
)

// TestGeminiCompleteKeyNotInError: a transport failure returns a *url.Error that
// embeds the full request URL, so the API key — sent only via the x-goog-api-key
// header — must never appear in the error string (key-leak gate).
func TestGeminiCompleteKeyNotInError(t *testing.T) {
	const key = "gk-secret-leak-canary"
	t.Setenv("GEMINI_API_KEY", key)
	g := Gemini{Model: "m", BaseURL: "http://127.0.0.1:1"}
	_, err := g.Complete("a", "b")
	if err == nil {
		t.Fatal("Complete = nil error, want request error")
	}
	if strings.Contains(err.Error(), key) {
		t.Fatalf("error string leaks the API key: %v", err)
	}
}
