//ff:func feature=llm type=helper control=sequence level=error
//ff:what TestFromFlagGeminiCLI — geminicli backend가 *GeminiCLI로 매핑되고(HTTP gemini:의 Gemini와 별개 case), model 토큰이 운반되며, geminicli:default가 빈 Model(CLI 기본)로 풀리는지 검증. REINS_GEMINI_SESSION을 비워 Stateless 기본을 고정.

package llm

import (
	"testing"
)

// TestFromFlagGeminiCLI: geminicli maps to *GeminiCLI (a separate case from the HTTP
// gemini:), carries the model token, and geminicli:default yields an empty Model.
func TestFromFlagGeminiCLI(t *testing.T) {
	t.Setenv("REINS_GEMINI_SESSION", "")

	b, err := FromFlag("geminicli:gemini-2.5-flash")
	if err != nil {
		t.Fatalf("FromFlag error: %v", err)
	}
	g, ok := b.(*GeminiCLI)
	if !ok {
		t.Fatalf("backend = %T, want *GeminiCLI", b)
	}
	if g.Model != "gemini-2.5-flash" {
		t.Fatalf("Model = %q, want gemini-2.5-flash", g.Model)
	}

	bd, err := FromFlag("geminicli:default")
	if err != nil {
		t.Fatalf("FromFlag(default) error: %v", err)
	}
	gd := bd.(*GeminiCLI)
	if gd.Model != "" {
		t.Fatalf("default Model = %q, want \"\"", gd.Model)
	}
}
