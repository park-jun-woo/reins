//ff:func feature=llm type=adapter control=sequence
//ff:what TestGeminiCompleteDefaultsVerbatim — 옵션 미지정 시 generationConfig가 현행과 동일(byte 수준): "maxOutputTokens":2048·"temperature":0이 그대로 실리는지 검증(후방호환 회귀 0).

package llm

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestGeminiCompleteDefaultsVerbatim: with no options the generationConfig carries
// maxOutputTokens 2048 and temperature 0.
func TestGeminiCompleteDefaultsVerbatim(t *testing.T) {
	t.Setenv("GEMINI_API_KEY", "gk")
	var rawBody string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, _ := io.ReadAll(r.Body)
		rawBody = string(data)
		io.WriteString(w, `{"candidates":[{"content":{"parts":[{"text":"x"}]}}]}`)
	}))
	defer srv.Close()

	g := Gemini{Model: "gemini-1.5-pro", BaseURL: srv.URL}
	if _, err := g.Complete("sys", "usr"); err != nil {
		t.Fatalf("Complete error: %v", err)
	}
	if !strings.Contains(rawBody, `"maxOutputTokens":2048`) {
		t.Fatalf("body %q missing \"maxOutputTokens\":2048", rawBody)
	}
	if !strings.Contains(rawBody, `"temperature":0`) {
		t.Fatalf("body %q missing \"temperature\":0", rawBody)
	}
}
