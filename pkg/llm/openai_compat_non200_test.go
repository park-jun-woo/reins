//ff:func feature=llm type=adapter control=sequence level=error
//ff:what TestOpenAICompatNon200 — 비-200 상태가 응답 body를 에러로 표면화하는지 검증.

package llm

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestOpenAICompatNon200: a non-200 status surfaces the body.
func TestOpenAICompatNon200(t *testing.T) {
	t.Setenv("XAI_API_KEY", "tok")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "bad-request")
	}))
	defer srv.Close()

	o := OpenAICompat{URL: srv.URL, Backend: "xai", Model: "grok"}
	_, err := o.Complete("a", "b")
	if err == nil || !strings.Contains(err.Error(), "bad-request") {
		t.Fatalf("error = %v, want body 'bad-request'", err)
	}
}
