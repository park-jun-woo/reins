//ff:func feature=llm type=adapter control=sequence level=error
//ff:what TestOpenAICompatBadURL — 제어문자가 포함된 URL이 http.NewRequest에서 실패하는지 검증.

package llm

import (
	"testing"
)

// TestOpenAICompatBadURL: a URL with a control character fails at http.NewRequest.
func TestOpenAICompatBadURL(t *testing.T) {
	t.Setenv("XAI_API_KEY", "tok")
	o := OpenAICompat{URL: "http://exa\x00mple", Backend: "xai", Model: "grok"}
	if _, err := o.Complete("a", "b"); err == nil {
		t.Fatal("Complete = nil error, want NewRequest error")
	}
}
