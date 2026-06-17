//ff:func feature=llm type=adapter control=sequence
//ff:what TestOllamaCompleteDefaultsVerbatim — 옵션 미지정 시 요청 body가 현행과 동일(byte 수준)함을 검증: options에 "num_predict":2048·"temperature":0이 그대로 실리고 "think" 키는 없어야 한다(후방호환 회귀 0).

package llm

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestOllamaCompleteDefaultsVerbatim: with no options set the request body still
// carries num_predict 2048 and temperature 0, and omits the think key.
func TestOllamaCompleteDefaultsVerbatim(t *testing.T) {
	var rawBody string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, _ := io.ReadAll(r.Body)
		rawBody = string(data)
		io.WriteString(w, `{"message":{"content":"x"}}`)
	}))
	defer srv.Close()

	o := Ollama{Model: "m", BaseURL: srv.URL}
	if _, err := o.Complete("sys", "usr"); err != nil {
		t.Fatalf("Complete error: %v", err)
	}
	if !strings.Contains(rawBody, `"num_predict":2048`) {
		t.Fatalf("body %q missing \"num_predict\":2048", rawBody)
	}
	if !strings.Contains(rawBody, `"temperature":0`) {
		t.Fatalf("body %q missing \"temperature\":0", rawBody)
	}
	if strings.Contains(rawBody, "think") {
		t.Fatalf("body %q must not contain think when Think is nil", rawBody)
	}
}
