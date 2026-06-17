//ff:func feature=llm type=adapter control=sequence
//ff:what TestOpenAICompatOptions — MaxOutputTokens=4096·Temperature=0.3 지정 시 요청 body의 max_tokens·temperature에 반영되는지 검증.

package llm

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestOpenAICompatOptions: MaxOutputTokens and Temperature map to max_tokens and
// temperature.
func TestOpenAICompatOptions(t *testing.T) {
	t.Setenv("XAI_API_KEY", "tok")
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(data, &gotBody)
		io.WriteString(w, `{"choices":[{"message":{"content":"x"}}]}`)
	}))
	defer srv.Close()

	temp := 0.3
	o := OpenAICompat{URL: srv.URL, Backend: "xai", Model: "grok", MaxOutputTokens: 4096, Temperature: &temp}
	if _, err := o.Complete("a", "b"); err != nil {
		t.Fatalf("Complete error: %v", err)
	}
	if mt := gotBody["max_tokens"].(float64); mt != 4096 {
		t.Fatalf("max_tokens = %v, want 4096", mt)
	}
	if tp := gotBody["temperature"].(float64); tp != 0.3 {
		t.Fatalf("temperature = %v, want 0.3", tp)
	}
}
