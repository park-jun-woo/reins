//ff:func feature=llm type=adapter control=sequence level=error
//ff:what TestOpenAICompatParseError — 잘못된 JSON 응답이 파싱 에러인지 검증.

package llm

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestOpenAICompatParseError: malformed JSON is a parse error.
func TestOpenAICompatParseError(t *testing.T) {
	t.Setenv("XAI_API_KEY", "tok")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "not json")
	}))
	defer srv.Close()

	o := OpenAICompat{URL: srv.URL, Backend: "xai", Model: "grok"}
	if _, err := o.Complete("a", "b"); err == nil {
		t.Fatal("Complete = nil error, want parse error")
	}
}
