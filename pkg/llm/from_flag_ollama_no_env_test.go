//ff:func feature=llm type=helper control=sequence level=error
//ff:what TestFromFlagOllamaNoEnv — REINS_OLLAMA_URL 미설정 시 BaseURL이 빈 문자열(어댑터가 호출 시점에 localhost로 폴백)인지 검증.

package llm

import (
	"testing"
)

// TestFromFlagOllamaNoEnv: without REINS_OLLAMA_URL the BaseURL is empty (adapter
// falls back to localhost at call time).
func TestFromFlagOllamaNoEnv(t *testing.T) {
	t.Setenv("REINS_OLLAMA_URL", "")
	b, err := FromFlag("ollama:llama3")
	if err != nil {
		t.Fatalf("FromFlag error: %v", err)
	}
	o := b.(Ollama)
	if o.BaseURL != "" {
		t.Fatalf("BaseURL = %q, want empty", o.BaseURL)
	}
}
