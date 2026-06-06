//ff:func feature=llm type=adapter control=sequence level=error
//ff:what TestGeminiCompleteReadError — 잘린 응답 body(Content-Length보다 적게 보내고 연결 끊음)가 io.ReadAll 실패를 일으키는지 검증.

package llm

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestGeminiCompleteReadError: a truncated response body makes io.ReadAll fail.
func TestGeminiCompleteReadError(t *testing.T) {
	t.Setenv("GEMINI_API_KEY", "gk")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "1000")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("short"))
		if hj, ok := w.(http.Hijacker); ok {
			conn, _, _ := hj.Hijack()
			conn.Close()
		}
	}))
	defer srv.Close()

	g := Gemini{Model: "m", BaseURL: srv.URL}
	if _, err := g.Complete("a", "b"); err == nil {
		t.Fatal("Complete = nil error, want read error")
	}
}
