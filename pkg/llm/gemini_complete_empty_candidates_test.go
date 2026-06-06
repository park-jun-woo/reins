//ff:func feature=llm type=adapter control=sequence level=error
//ff:what TestGeminiCompleteEmptyCandidates — candidates가 0개(또는 빈 parts)면 에러인지 검증.

package llm

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestGeminiCompleteEmptyCandidates: zero candidates (or empty parts) is an error.
func TestGeminiCompleteEmptyCandidates(t *testing.T) {
	t.Setenv("GEMINI_API_KEY", "gk")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"candidates":[]}`)
	}))
	defer srv.Close()

	g := Gemini{Model: "m", BaseURL: srv.URL}
	if _, err := g.Complete("a", "b"); err == nil {
		t.Fatal("Complete = nil error, want empty-candidates error")
	}
}
