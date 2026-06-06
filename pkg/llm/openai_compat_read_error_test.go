//ff:func feature=llm type=adapter control=sequence level=error
//ff:what TestOpenAICompatReadError — 광고한 길이보다 짧게 보내고 연결을 끊는 서버가 응답 body의 io.ReadAll 실패를 일으키는지 검증.

package llm

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestOpenAICompatReadError: a server that advertises more bytes than it sends and
// then hangs up makes io.ReadAll of the response body fail.
func TestOpenAICompatReadError(t *testing.T) {
	t.Setenv("XAI_API_KEY", "tok")
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

	o := OpenAICompat{URL: srv.URL, Backend: "xai", Model: "grok"}
	if _, err := o.Complete("a", "b"); err == nil {
		t.Fatal("Complete = nil error, want read error")
	}
}
