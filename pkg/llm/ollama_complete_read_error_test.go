//ff:func feature=llm type=adapter control=sequence level=error
//ff:what TestOllamaCompleteReadError — 광고한 길이보다 짧게 보내고 연결을 끊는 응답 body가 io.ReadAll 실패를 일으키는지 검증.

package llm

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestOllamaCompleteReadError: a truncated response body (advertised longer than
// sent, then the connection drops) makes io.ReadAll of the body fail.
func TestOllamaCompleteReadError(t *testing.T) {
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

	o := Ollama{Model: "m", BaseURL: srv.URL}
	if _, err := o.Complete("a", "b"); err == nil {
		t.Fatal("Complete = nil error, want read error")
	}
}
