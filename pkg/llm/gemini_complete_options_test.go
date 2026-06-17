//ff:func feature=llm type=adapter control=sequence
//ff:what TestGeminiCompleteOptions — MaxOutputTokens=4096·Temperature=0.9 지정 시 generationConfig.maxOutputTokens·temperature에 반영되는지 검증.

package llm

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestGeminiCompleteOptions: MaxOutputTokens and Temperature map into generationConfig.
func TestGeminiCompleteOptions(t *testing.T) {
	t.Setenv("GEMINI_API_KEY", "gk")
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(data, &gotBody)
		io.WriteString(w, `{"candidates":[{"content":{"parts":[{"text":"x"}]}}]}`)
	}))
	defer srv.Close()

	temp := 0.9
	g := Gemini{Model: "gemini-1.5-pro", BaseURL: srv.URL, MaxOutputTokens: 4096, Temperature: &temp}
	if _, err := g.Complete("a", "b"); err != nil {
		t.Fatalf("Complete error: %v", err)
	}
	cfg := gotBody["generationConfig"].(map[string]any)
	if mt := cfg["maxOutputTokens"].(float64); mt != 4096 {
		t.Fatalf("maxOutputTokens = %v, want 4096", mt)
	}
	if tp := cfg["temperature"].(float64); tp != 0.9 {
		t.Fatalf("temperature = %v, want 0.9", tp)
	}
}
