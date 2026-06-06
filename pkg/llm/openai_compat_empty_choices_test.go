//ff:func feature=llm type=adapter control=sequence level=error
//ff:what TestOpenAICompatEmptyChoices — choices가 0개면 에러인지 검증.

package llm

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestOpenAICompatEmptyChoices: zero choices is an error.
func TestOpenAICompatEmptyChoices(t *testing.T) {
	t.Setenv("XAI_API_KEY", "tok")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"choices":[]}`)
	}))
	defer srv.Close()

	o := OpenAICompat{URL: srv.URL, Backend: "xai", Model: "grok"}
	if _, err := o.Complete("a", "b"); err == nil {
		t.Fatal("Complete = nil error, want empty-choices error")
	}
}
