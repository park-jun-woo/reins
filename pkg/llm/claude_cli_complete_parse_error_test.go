//ff:func feature=llm type=adapter control=sequence level=error
//ff:what TestClaudeCLICompleteParseError — execClaude가 비JSON stdout을 반환하면 Complete가 파싱 에러를 내는지 무서브프로세스 스텁으로 검증.

package llm

import (
	"context"
	"testing"
)

// TestClaudeCLICompleteParseError: malformed JSON stdout is a parse error.
func TestClaudeCLICompleteParseError(t *testing.T) {
	orig := execClaude
	defer func() { execClaude = orig }()

	execClaude = func(ctx context.Context, bin string, argv []string, stdin string) (string, string, error) {
		return "not json", "", nil
	}

	c := &ClaudeCLI{}
	if _, err := c.Complete("SYS", "USR"); err == nil {
		t.Fatal("Complete = nil error, want parse error")
	}
}
