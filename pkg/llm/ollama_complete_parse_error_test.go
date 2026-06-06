//ff:func feature=llm type=adapter control=sequence level=error
//ff:what TestOllamaCompleteParseError — 잘못된 JSON 응답이 파싱 에러인지 검증.

package llm

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestOllamaCompleteParseError: malformed JSON yields a parse error.
func TestOllamaCompleteParseError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "not json")
	}))
	defer srv.Close()

	o := Ollama{Model: "m", BaseURL: srv.URL}
	if _, err := o.Complete("a", "b"); err == nil {
		t.Fatal("Complete = nil error, want parse error")
	}
}
