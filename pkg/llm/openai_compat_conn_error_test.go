//ff:func feature=llm type=adapter control=sequence level=error
//ff:what TestOpenAICompatConnError — 도달 불가 URL이 요청 에러인지 검증.

package llm

import (
	"testing"
)

// TestOpenAICompatConnError: an unreachable URL is a request error.
func TestOpenAICompatConnError(t *testing.T) {
	t.Setenv("XAI_API_KEY", "tok")
	o := OpenAICompat{URL: "http://127.0.0.1:1", Backend: "xai", Model: "grok"}
	if _, err := o.Complete("a", "b"); err == nil {
		t.Fatal("Complete = nil error, want request error")
	}
}
