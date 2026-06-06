//ff:func feature=llm type=adapter control=sequence level=error
//ff:what TestDoOllamaRequestMarshalError — JSON 직렬화 불가 body(채널)가 HTTP 호출 전에 마샬 에러로 실패하는지 검증.

package llm

import (
	"testing"
)

// TestDoOllamaRequestMarshalError: an unmarshalable body fails before any HTTP call.
func TestDoOllamaRequestMarshalError(t *testing.T) {
	// channels can't be JSON-marshaled.
	if _, err := doOllamaRequest("http://unused", make(chan int)); err == nil {
		t.Fatal("doOllamaRequest = nil error, want marshal error")
	}
}
