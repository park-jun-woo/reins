//ff:func feature=llm type=helper control=sequence level=error
//ff:what TestFromFlagOllamaQuery — "ollama:gemma4:e4b?max_output_tokens=8192&num_ctx=16384&temperature=0.5&think=false"가 콜론 포함 model을 보존하면서 Ollama 필드(MaxOutputTokens·NumCtx·Temperature·Think)로 매핑되는지 검증.

package llm

import "testing"

// TestFromFlagOllamaQuery: a query maps onto Ollama fields while the colon-bearing
// model name is preserved.
func TestFromFlagOllamaQuery(t *testing.T) {
	t.Setenv("REINS_OLLAMA_URL", "")
	b, err := FromFlag("ollama:gemma4:e4b?max_output_tokens=8192&num_ctx=16384&temperature=0.5&think=false")
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
	if o.MaxOutputTokens != 8192 {
		t.Fatalf("MaxOutputTokens = %d, want 8192", o.MaxOutputTokens)
	}
	if o.NumCtx != 16384 {
		t.Fatalf("NumCtx = %d, want 16384", o.NumCtx)
	}
	if o.Temperature == nil || *o.Temperature != 0.5 {
		t.Fatalf("Temperature = %v, want 0.5", o.Temperature)
	}
	if o.Think == nil || *o.Think != false {
		t.Fatalf("Think = %v, want false", o.Think)
	}
}
