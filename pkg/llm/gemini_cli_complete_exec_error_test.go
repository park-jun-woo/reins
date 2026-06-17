//ff:func feature=llm type=adapter control=sequence level=error
//ff:what TestGeminiCLICompleteExecError — execGemini가 비정상 종료(non-nil err) 시 그 에러가 stderr를 동봉해 래핑되어 Complete에서 표면화되는지 검증.

package llm

import (
	"context"
	"errors"
	"strings"
	"testing"
)

// TestGeminiCLICompleteExecError: a non-nil exec error surfaces with stderr in the
// wrapped error.
func TestGeminiCLICompleteExecError(t *testing.T) {
	orig := execGemini
	defer func() { execGemini = orig }()

	execGemini = func(ctx context.Context, bin string, argv []string, stdin string) (string, string, error) {
		return "", "boom-stderr", errors.New("exit 1")
	}

	c := &GeminiCLI{}
	_, err := c.Complete("SYS", "USR")
	if err == nil {
		t.Fatal("Complete = nil error, want exec error")
	}
	if !strings.Contains(err.Error(), "boom-stderr") {
		t.Fatalf("error = %q, want it to contain stderr", err.Error())
	}
}
