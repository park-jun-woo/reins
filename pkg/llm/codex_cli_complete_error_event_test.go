//ff:func feature=llm type=adapter control=sequence level=error
//ff:what TestCodexCLICompleteErrorEvent — JSONL 스트림에 `{"type":"error",...}` 이벤트가 오면 Complete가 그 이벤트를 에러로 표면화하는지(agent_message가 뒤따라도 에러 우선) 검증.

package llm

import (
	"context"
	"strings"
	"testing"
)

// TestCodexCLICompleteErrorEvent: an `error` event in the JSONL stream surfaces as a
// Complete error (it short-circuits even if an agent_message would follow).
func TestCodexCLICompleteErrorEvent(t *testing.T) {
	orig := execCodex
	defer func() { execCodex = orig }()

	execCodex = func(ctx context.Context, bin string, argv []string, stdin string) (string, string, error) {
		return `{"type":"thread.started","thread_id":"T1"}
{"type":"error","message":"model overloaded"}`, "", nil
	}

	c := &CodexCLI{}
	_, err := c.Complete("SYS", "USR")
	if err == nil {
		t.Fatal("Complete = nil error, want error-event error")
	}
	if !strings.Contains(err.Error(), "error event") {
		t.Fatalf("error = %q, want it to mention error event", err.Error())
	}
}
