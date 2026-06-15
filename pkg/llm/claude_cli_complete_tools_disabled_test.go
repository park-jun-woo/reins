//ff:func feature=llm type=adapter control=sequence level=error
//ff:what TestClaudeCLICompleteToolsDisabled — argv에 --tools가 빈 값("")으로 실려 전 도구가 비활성화됨을 execClaude 스텁으로 검증. 도구 비활성이 없으면 모델이 도구를 시도→--max-turns 1 소진→error_max_turns(stop_reason tool_use)로 생성이 하드 에러가 되어 루프를 중단시키므로, L0 순수 텍스트 생성의 불변식이다.

package llm

import (
	"context"
	"testing"
)

// TestClaudeCLICompleteToolsDisabled: every Complete passes --tools "" so the model
// cannot attempt a tool call. Without it, a tool attempt under --max-turns 1 exits
// error_max_turns before producing a result — a regression that turns generation
// into a loop-aborting backend error.
func TestClaudeCLICompleteToolsDisabled(t *testing.T) {
	orig := execClaude
	defer func() { execClaude = orig }()

	var gotArgv []string
	execClaude = func(ctx context.Context, bin string, argv []string, stdin string) (string, string, error) {
		gotArgv = argv
		return `{"result":"OK","session_id":"S1","is_error":false}`, "", nil
	}

	c := &ClaudeCLI{}
	if _, err := c.Complete("SYS", "USR"); err != nil {
		t.Fatalf("Complete error: %v", err)
	}
	if !hasFlag(gotArgv, "--tools") {
		t.Fatalf("argv %v missing --tools", gotArgv)
	}
	if v := flagValue(gotArgv, "--tools"); v != "" {
		t.Fatalf("--tools = %q, want empty (all tools disabled)", v)
	}
}
