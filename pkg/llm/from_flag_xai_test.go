//ff:func feature=llm type=helper control=sequence level=error
//ff:what TestFromFlagXai — xai backend가 x.ai endpoint의 OpenAICompat(Backend·Model·URL)로 매핑되는지 검증.

package llm

import (
	"testing"
)

// TestFromFlagXai: xai maps to OpenAICompat at the x.ai endpoint.
func TestFromFlagXai(t *testing.T) {
	b, err := FromFlag("xai:grok-2")
	if err != nil {
		t.Fatalf("FromFlag error: %v", err)
	}
	o, ok := b.(OpenAICompat)
	if !ok {
		t.Fatalf("backend = %T, want OpenAICompat", b)
	}
	if o.Model != "grok-2" {
		t.Fatalf("Model = %q", o.Model)
	}
	if o.Backend != "xai" {
		t.Fatalf("Backend = %q, want xai", o.Backend)
	}
	if o.URL != "https://api.x.ai/v1/chat/completions" {
		t.Fatalf("URL = %q", o.URL)
	}
}
