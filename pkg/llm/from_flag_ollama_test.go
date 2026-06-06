//ff:func feature=llm type=helper control=sequence level=error
//ff:what TestFromFlagOllama — ollama backend, model에 ':' 허용(gemma4:e4b), BaseURL이 env REINS_OLLAMA_URL 오버라이드를 받는지 검증.

package llm

import (
	"testing"
)

// TestFromFlagOllama: ollama backend, model may contain colons, BaseURL picks up
// REINS_OLLAMA_URL.
func TestFromFlagOllama(t *testing.T) {
	t.Setenv("REINS_OLLAMA_URL", "http://example:1234")
	b, err := FromFlag("ollama:gemma4:e4b")
	if err != nil {
		t.Fatalf("FromFlag error: %v", err)
	}
	o, ok := b.(Ollama)
	if !ok {
		t.Fatalf("backend = %T, want Ollama", b)
	}
	if o.Model != "gemma4:e4b" {
		t.Fatalf("Model = %q, want gemma4:e4b", o.Model)
	}
	if o.BaseURL != "http://example:1234" {
		t.Fatalf("BaseURL = %q, want override", o.BaseURL)
	}
}
