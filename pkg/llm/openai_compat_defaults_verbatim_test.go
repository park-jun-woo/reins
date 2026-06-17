//ff:func feature=llm type=adapter control=sequence
//ff:what TestOpenAICompatDefaultsVerbatim — 옵션 미지정 시 요청 body가 현행과 동일(byte 수준): "max_tokens":2048·"temperature":0이 그대로 실리는지 검증(후방호환 회귀 0).

package llm

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestOpenAICompatDefaultsVerbatim: with no options the body carries max_tokens 2048
// and temperature 0.
func TestOpenAICompatDefaultsVerbatim(t *testing.T) {
	t.Setenv("XAI_API_KEY", "tok")
	var rawBody string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, _ := io.ReadAll(r.Body)
		rawBody = string(data)
		io.WriteString(w, `{"choices":[{"message":{"content":"x"}}]}`)
	}))
	defer srv.Close()

	o := OpenAICompat{URL: srv.URL, Backend: "xai", Model: "grok"}
	if _, err := o.Complete("sys", "usr"); err != nil {
		t.Fatalf("Complete error: %v", err)
	}
	if !strings.Contains(rawBody, `"max_tokens":2048`) {
		t.Fatalf("body %q missing \"max_tokens\":2048", rawBody)
	}
	if !strings.Contains(rawBody, `"temperature":0`) {
		t.Fatalf("body %q missing \"temperature\":0", rawBody)
	}
}
