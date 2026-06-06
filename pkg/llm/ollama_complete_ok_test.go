//ff:func feature=llm type=adapter control=sequence
//ff:what TestOllamaCompleteOK — 200 응답의 message.content를 파싱·반환하고 요청 body가 model·system+user 메시지·자동 num_ctx를 담는지 httptest로 검증.

package llm

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestOllamaCompleteOK: a 200 response with message.content is parsed and returned,
// and the request body carries model, system+user messages, and the auto num_ctx.
func TestOllamaCompleteOK(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/chat" {
			t.Errorf("path = %q, want /api/chat", r.URL.Path)
		}
		data, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(data, &gotBody)
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{"message":{"content":"hello world"}}`)
	}))
	defer srv.Close()

	o := Ollama{Model: "m", BaseURL: srv.URL}
	got, err := o.Complete("sys", "usr")
	if err != nil {
		t.Fatalf("Complete error: %v", err)
	}
	if got != "hello world" {
		t.Fatalf("content = %q, want hello world", got)
	}
	if gotBody["model"] != "m" {
		t.Fatalf("model = %v", gotBody["model"])
	}
	opts, _ := gotBody["options"].(map[string]any)
	if opts == nil {
		t.Fatalf("options missing: %v", gotBody)
	}
	// num_ctx is auto-sized (NumCtx==0) and must be a valid context step.
	if nc, _ := opts["num_ctx"].(float64); nc < 2048 || nc > 32768 {
		t.Fatalf("num_ctx = %v out of range", opts["num_ctx"])
	}
}
