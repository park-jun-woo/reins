//ff:func feature=llm type=adapter control=sequence level=error
//ff:what TestClaudeCLICompleteExecError — execClaude가 non-nil err+stderr를 반환하면 Complete가 stderr를 담은 래핑 에러를 내는지 무서브프로세스 스텁으로 검증.

package llm

import (
	"context"
	"errors"
	"strings"
	"testing"
)

// TestClaudeCLICompleteExecError: a non-nil exec error surfaces with stderr in
// the wrapped error.
func TestClaudeCLICompleteExecError(t *testing.T) {
	orig := execClaude
	defer func() { execClaude = orig }()

	execClaude = func(ctx context.Context, bin string, argv []string, stdin string) (string, string, error) {
		return "", "boom-stderr", errors.New("exit 1")
	}

	c := &ClaudeCLI{}
	_, err := c.Complete("SYS", "USR")
	if err == nil {
		t.Fatal("Complete = nil error, want exec error")
	}
	if !strings.Contains(err.Error(), "boom-stderr") {
		t.Fatalf("error = %q, want it to contain stderr", err.Error())
	}
}
