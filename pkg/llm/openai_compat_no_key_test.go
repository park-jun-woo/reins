//ff:func feature=llm type=adapter control=sequence level=error
//ff:what TestOpenAICompatNoKey — API 키 미설정 시 HTTP 호출 전에 에러로 단락되는지 검증.

package llm

import (
	"testing"
)

// TestOpenAICompatNoKey: a missing API key errors before any HTTP call.
func TestOpenAICompatNoKey(t *testing.T) {
	t.Setenv("XAI_API_KEY", "")
	o := OpenAICompat{URL: "http://unused", Backend: "xai", Model: "grok"}
	if _, err := o.Complete("a", "b"); err == nil {
		t.Fatal("Complete = nil error, want missing-key error")
	}
}
