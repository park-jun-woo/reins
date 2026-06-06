//ff:func feature=llm type=adapter control=sequence level=error
//ff:what TestOllamaCompleteConnError — 도달 불가 서버면 요청 에러인지 검증.

package llm

import (
	"testing"
)

// TestOllamaCompleteConnError: an unreachable server yields a request error.
func TestOllamaCompleteConnError(t *testing.T) {
	o := Ollama{Model: "m", BaseURL: "http://127.0.0.1:1"}
	if _, err := o.Complete("a", "b"); err == nil {
		t.Fatal("Complete = nil error, want request error")
	}
}
