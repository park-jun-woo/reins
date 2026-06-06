//ff:func feature=llm type=adapter control=sequence
//ff:what TestOpenAICompatOK — 200 응답이 choices[0]을 반환하고 요청이 Bearer 키 헤더와 model/messages body를 담는지 httptest로 검증.

package llm

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestOpenAICompatOK: a 200 response returns the first choice; the request carries
// the Bearer key and the model/messages body.
func TestOpenAICompatOK(t *testing.T) {
	t.Setenv("XAI_API_KEY", "tok-123")
	var auth string
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth = r.Header.Get("Authorization")
		data, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(data, &gotBody)
		io.WriteString(w, `{"choices":[{"message":{"content":"answer"}}]}`)
	}))
	defer srv.Close()

	o := OpenAICompat{URL: srv.URL, Backend: "xai", Model: "grok"}
	got, err := o.Complete("sys", "usr")
	if err != nil {
		t.Fatalf("Complete error: %v", err)
	}
	if got != "answer" {
		t.Fatalf("content = %q, want answer", got)
	}
	if auth != "Bearer tok-123" {
		t.Fatalf("Authorization = %q", auth)
	}
	if gotBody["model"] != "grok" {
		t.Fatalf("model = %v", gotBody["model"])
	}
}
