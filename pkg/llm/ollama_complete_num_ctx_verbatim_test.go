//ff:func feature=llm type=adapter control=sequence
//ff:what TestOllamaCompleteNumCtxVerbatim — NumCtx가 0이 아니면 그 값이 그대로 요청 body의 num_ctx에 쓰이는지 검증.

package llm

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestOllamaCompleteNumCtxVerbatim: a non-zero NumCtx is used as-is.
func TestOllamaCompleteNumCtxVerbatim(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(data, &gotBody)
		io.WriteString(w, `{"message":{"content":"x"}}`)
	}))
	defer srv.Close()

	o := Ollama{Model: "m", BaseURL: srv.URL, NumCtx: 65536}
	if _, err := o.Complete("a", "b"); err != nil {
		t.Fatalf("Complete error: %v", err)
	}
	opts := gotBody["options"].(map[string]any)
	if nc := opts["num_ctx"].(float64); nc != 65536 {
		t.Fatalf("num_ctx = %v, want 65536 verbatim", nc)
	}
}
