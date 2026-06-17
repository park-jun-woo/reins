//ff:func feature=llm type=helper control=sequence level=error
//ff:what TestFromFlagXaiQuery — "xai:grok?max_output_tokens=4096&temperature=0.2"가 OpenAICompat의 MaxOutputTokens·Temperature로 매핑되는지(허용 키만) 검증.

package llm

import "testing"

// TestFromFlagXaiQuery: xai accepts max_output_tokens and temperature.
func TestFromFlagXaiQuery(t *testing.T) {
	b, err := FromFlag("xai:grok?max_output_tokens=4096&temperature=0.2")
	if err != nil {
		t.Fatalf("FromFlag error: %v", err)
	}
	o, ok := b.(OpenAICompat)
	if !ok {
		t.Fatalf("backend = %T, want OpenAICompat", b)
	}
	if o.MaxOutputTokens != 4096 {
		t.Fatalf("MaxOutputTokens = %d, want 4096", o.MaxOutputTokens)
	}
	if o.Temperature == nil || *o.Temperature != 0.2 {
		t.Fatalf("Temperature = %v, want 0.2", o.Temperature)
	}
}
