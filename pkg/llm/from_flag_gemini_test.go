//ff:func feature=llm type=helper control=sequence level=error
//ff:what TestFromFlagGemini — gemini backend가 Gemini(Model)로 매핑되는지 검증.

package llm

import (
	"testing"
)

// TestFromFlagGemini: gemini maps to Gemini.
func TestFromFlagGemini(t *testing.T) {
	b, err := FromFlag("gemini:gemini-1.5-pro")
	if err != nil {
		t.Fatalf("FromFlag error: %v", err)
	}
	g, ok := b.(Gemini)
	if !ok {
		t.Fatalf("backend = %T, want Gemini", b)
	}
	if g.Model != "gemini-1.5-pro" {
		t.Fatalf("Model = %q", g.Model)
	}
}
