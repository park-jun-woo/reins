//ff:func feature=llm type=adapter control=sequence
//ff:what TestOllamaCompleteTemperature — Temperature 포인터가 설정되면 options.temperature에 그 값이 반영되는지 검증(nil은 별도 verbatim 테스트가 0 고정 확인).

package llm

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestOllamaCompleteTemperature: a non-nil Temperature is reflected in options.
func TestOllamaCompleteTemperature(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(data, &gotBody)
		io.WriteString(w, `{"message":{"content":"x"}}`)
	}))
	defer srv.Close()

	temp := 0.7
	o := Ollama{Model: "m", BaseURL: srv.URL, Temperature: &temp}
	if _, err := o.Complete("a", "b"); err != nil {
		t.Fatalf("Complete error: %v", err)
	}
	opts := gotBody["options"].(map[string]any)
	if got := opts["temperature"].(float64); got != 0.7 {
		t.Fatalf("temperature = %v, want 0.7", got)
	}
}
