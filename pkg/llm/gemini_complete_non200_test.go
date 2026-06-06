//ff:func feature=llm type=adapter control=sequence level=error
//ff:what TestGeminiCompleteNon200 — 비-200 상태가 응답 body를 에러로 표면화하는지 검증.

package llm

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestGeminiCompleteNon200: a non-200 status surfaces the body.
func TestGeminiCompleteNon200(t *testing.T) {
	t.Setenv("GEMINI_API_KEY", "gk")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		io.WriteString(w, "denied")
	}))
	defer srv.Close()

	g := Gemini{Model: "m", BaseURL: srv.URL}
	_, err := g.Complete("a", "b")
	if err == nil || !strings.Contains(err.Error(), "denied") {
		t.Fatalf("error = %v, want body 'denied'", err)
	}
}
