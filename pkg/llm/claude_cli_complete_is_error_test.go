//ff:func feature=llm type=adapter control=sequence level=error
//ff:what TestClaudeCLICompleteIsError — is_error:true+subtype 응답이면 Complete가 subtype을 담은 에러를 내는지 무서브프로세스 스텁으로 검증.

package llm

import (
	"context"
	"strings"
	"testing"
)

// TestClaudeCLICompleteIsError: is_error:true surfaces an error carrying subtype.
func TestClaudeCLICompleteIsError(t *testing.T) {
	orig := execClaude
	defer func() { execClaude = orig }()

	execClaude = func(ctx context.Context, bin string, argv []string, stdin string) (string, string, error) {
		return `{"result":"","is_error":true,"subtype":"max_turns"}`, "", nil
	}

	c := &ClaudeCLI{}
	_, err := c.Complete("SYS", "USR")
	if err == nil {
		t.Fatal("Complete = nil error, want is_error")
	}
	if !strings.Contains(err.Error(), "max_turns") {
		t.Fatalf("error = %q, want it to contain subtype", err.Error())
	}
}
