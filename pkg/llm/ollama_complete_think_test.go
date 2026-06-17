//ff:func feature=llm type=adapter control=sequence
//ff:what TestOllamaCompleteThink — Think=false 지정 시 options.think이 false로 실리는지 검증(nil은 think 키 부재 — verbatim 테스트가 확인).

package llm

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestOllamaCompleteThink: Think=false sets options.think to false.
func TestOllamaCompleteThink(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(data, &gotBody)
		io.WriteString(w, `{"message":{"content":"x"}}`)
	}))
	defer srv.Close()

	think := false
	o := Ollama{Model: "m", BaseURL: srv.URL, Think: &think}
	if _, err := o.Complete("a", "b"); err != nil {
		t.Fatalf("Complete error: %v", err)
	}
	opts := gotBody["options"].(map[string]any)
	got, ok := opts["think"]
	if !ok {
		t.Fatalf("options missing think key: %v", opts)
	}
	if got.(bool) != false {
		t.Fatalf("think = %v, want false", got)
	}
}
