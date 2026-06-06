//ff:func feature=llm type=adapter control=sequence level=error
//ff:what TestGeminiCompleteParseError — 잘못된 JSON 응답이 파싱 에러인지 검증.

package llm

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestGeminiCompleteParseError: malformed JSON is a parse error.
func TestGeminiCompleteParseError(t *testing.T) {
	t.Setenv("GEMINI_API_KEY", "gk")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "not json")
	}))
	defer srv.Close()

	g := Gemini{Model: "m", BaseURL: srv.URL}
	if _, err := g.Complete("a", "b"); err == nil {
		t.Fatal("Complete = nil error, want parse error")
	}
}
