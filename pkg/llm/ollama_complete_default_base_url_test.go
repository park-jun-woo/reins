//ff:func feature=llm type=adapter control=sequence level=error
//ff:what TestOllamaCompleteDefaultBaseURL — 빈 BaseURL이 localhost 기본값으로 폴백하는 분기를 커버한다(테스트에 ollama 미실행이라 요청은 에러; 분기 자체가 커버 대상).

package llm

import (
	"testing"
)

// TestOllamaCompleteDefaultBaseURL: an empty BaseURL falls back to the localhost
// default (exercising the default branch). No ollama runs in tests, so the request
// errors — that is fine; the branch is what we cover.
func TestOllamaCompleteDefaultBaseURL(t *testing.T) {
	o := Ollama{Model: "m"} // BaseURL empty ⇒ defaultOllamaBaseURL
	if _, err := o.Complete("a", "b"); err == nil {
		t.Skip("an ollama server is unexpectedly reachable on localhost:11434")
	}
}
