//ff:func feature=llm type=adapter control=sequence level=error
//ff:what TestOllamaCompleteNon200 — 비-200 상태가 응답 body를 포함한 에러인지 검증.

package llm

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestOllamaCompleteNon200: a non-200 status yields an error including the body.
func TestOllamaCompleteNon200(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "boom")
	}))
	defer srv.Close()

	o := Ollama{Model: "m", BaseURL: srv.URL}
	_, err := o.Complete("a", "b")
	if err == nil {
		t.Fatal("Complete = nil error, want non-200 error")
	}
	if !strings.Contains(err.Error(), "boom") {
		t.Fatalf("error = %v, want body 'boom'", err)
	}
}
