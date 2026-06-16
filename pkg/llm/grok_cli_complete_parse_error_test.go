//ff:func feature=llm type=adapter control=sequence level=error
//ff:what TestGrokCLICompleteParseError — execGrok이 JSON이 아닌 stdout을 돌려줄 때 Complete가 파싱 에러를 표면화하는지 검증.

package llm

import (
	"context"
	"testing"
)

// TestGrokCLICompleteParseError: malformed JSON stdout is a parse error.
func TestGrokCLICompleteParseError(t *testing.T) {
	orig := execGrok
	defer func() { execGrok = orig }()

	execGrok = func(ctx context.Context, bin string, argv []string, stdin string) (string, string, error) {
		return "not json", "", nil
	}

	g := &GrokCLI{}
	if _, err := g.Complete("SYS", "USR"); err == nil {
		t.Fatal("Complete = nil error, want parse error")
	}
}
