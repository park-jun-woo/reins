//ff:func feature=llm type=adapter control=sequence level=error
//ff:what TestCodexCLICompleteParseError — JSONL 스트림에 agent_message 이벤트가 하나도 없으면(중간 이벤트만) Complete가 "no assistant message" 에러를 표면화하는지 검증.

package llm

import (
	"context"
	"strings"
	"testing"
)

// TestCodexCLICompleteParseError: a JSONL stream with no agent_message event (only
// intermediate events) surfaces a "no assistant message" error.
func TestCodexCLICompleteParseError(t *testing.T) {
	orig := execCodex
	defer func() { execCodex = orig }()

	execCodex = func(ctx context.Context, bin string, argv []string, stdin string) (string, string, error) {
		return `{"type":"thread.started","thread_id":"T1"}
{"type":"turn.started"}
{"type":"turn.completed"}`, "", nil
	}

	c := &CodexCLI{}
	_, err := c.Complete("SYS", "USR")
	if err == nil {
		t.Fatal("Complete = nil error, want parse error")
	}
	if !strings.Contains(err.Error(), "no assistant message") {
		t.Fatalf("error = %q, want it to mention no assistant message", err.Error())
	}
}
