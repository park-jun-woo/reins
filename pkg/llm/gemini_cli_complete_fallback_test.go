//ff:func feature=llm type=adapter control=sequence level=error
//ff:what TestGeminiCLICompleteFallback — gemini가 -o json을 미준수하고 plain 텍스트를 뱉는 알려진 이슈(#11184)에서 Complete가 폴백으로 trimmed plain stdout을 response로 반환하는지 execGemini 스텁으로 검증.

package llm

import (
	"context"
	"testing"
)

// TestGeminiCLICompleteFallback: when gemini ignores -o json and emits plain text
// (#11184), Complete falls back to the trimmed plain stdout as the response.
func TestGeminiCLICompleteFallback(t *testing.T) {
	orig := execGemini
	defer func() { execGemini = orig }()

	execGemini = func(ctx context.Context, bin string, argv []string, stdin string) (string, string, error) {
		return "  plain answer, not json\n", "", nil
	}

	c := &GeminiCLI{}
	got, err := c.Complete("SYS", "USR")
	if err != nil {
		t.Fatalf("Complete error: %v", err)
	}
	if got != "plain answer, not json" {
		t.Fatalf("result = %q, want trimmed plain text fallback", got)
	}
}
