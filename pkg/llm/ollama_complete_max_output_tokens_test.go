//ff:func feature=llm type=adapter control=sequence
//ff:what TestOllamaCompleteMaxOutputTokens — MaxOutputTokens=8192 지정 시 options.num_predict==8192이고, 같은 값이 autoNumCtx 예비분으로 스레딩되어 num_ctx가 기본(2048)보다 계단 상향(8192)되는지 검증(Phase017 A의 핵심 부수효과).

package llm

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestOllamaCompleteMaxOutputTokens: MaxOutputTokens feeds both num_predict and the
// autoNumCtx reserve, so num_ctx steps up from the default. Prompt "s"+"u" (2 ASCII
// bytes) ⇒ need = 0 + 8192 ⇒ step 8192; the default reserve would give 2048.
func TestOllamaCompleteMaxOutputTokens(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(data, &gotBody)
		io.WriteString(w, `{"message":{"content":"x"}}`)
	}))
	defer srv.Close()

	o := Ollama{Model: "m", BaseURL: srv.URL, MaxOutputTokens: 8192}
	if _, err := o.Complete("s", "u"); err != nil {
		t.Fatalf("Complete error: %v", err)
	}
	opts := gotBody["options"].(map[string]any)
	if np := opts["num_predict"].(float64); np != 8192 {
		t.Fatalf("num_predict = %v, want 8192", np)
	}
	if nc := opts["num_ctx"].(float64); nc != 8192 {
		t.Fatalf("num_ctx = %v, want 8192 (reserve threaded; default would be 2048)", nc)
	}
}
