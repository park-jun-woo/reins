//ff:func feature=llm type=adapter control=sequence level=error
//ff:what TestGeminiCLICompleteErrorEnvelope — json 봉투의 `error` 필드가 non-null이면 Complete가 그 봉투를 에러로 표면화하는지(response가 있어도 error 우선) 검증.

package llm

import (
	"context"
	"strings"
	"testing"
)

// TestGeminiCLICompleteErrorEnvelope: a non-null `error` field in the json envelope
// surfaces as a Complete error.
func TestGeminiCLICompleteErrorEnvelope(t *testing.T) {
	orig := execGemini
	defer func() { execGemini = orig }()

	execGemini = func(ctx context.Context, bin string, argv []string, stdin string) (string, string, error) {
		return `{"response":null,"error":{"message":"quota exceeded"}}`, "", nil
	}

	c := &GeminiCLI{}
	_, err := c.Complete("SYS", "USR")
	if err == nil {
		t.Fatal("Complete = nil error, want error-envelope error")
	}
	if !strings.Contains(err.Error(), "error envelope") {
		t.Fatalf("error = %q, want it to mention error envelope", err.Error())
	}
}
