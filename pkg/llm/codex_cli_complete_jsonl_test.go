//ff:func feature=llm type=adapter control=sequence level=error
//ff:what TestCodexCLICompleteJSONL — 다중 JSONL 이벤트(thread.started·중간 비-agent 이벤트·여러 agent_message·turn.completed) 중 **최종 agent_message** 텍스트만 결과로 취하고 중간 이벤트·미파싱 노이즈 줄을 무시하는지 execCodex 스텁으로 검증.

package llm

import (
	"context"
	"testing"
)

// TestCodexCLICompleteJSONL: from a multi-event JSONL stream the result is the LAST
// agent_message text; intermediate events and unparseable noise lines are ignored.
func TestCodexCLICompleteJSONL(t *testing.T) {
	orig := execCodex
	defer func() { execCodex = orig }()

	execCodex = func(ctx context.Context, bin string, argv []string, stdin string) (string, string, error) {
		return `{"type":"thread.started","thread_id":"T1"}
{"type":"turn.started"}
not-json-noise
{"type":"item.completed","item":{"type":"reasoning","text":"thinking"}}
{"type":"item.completed","item":{"type":"agent_message","text":"first"}}
{"type":"item.completed","item":{"type":"agent_message","text":"FINAL"}}
{"type":"turn.completed","usage":{"input_tokens":1}}`, "", nil
	}

	c := &CodexCLI{}
	got, err := c.Complete("SYS", "USR")
	if err != nil {
		t.Fatalf("Complete error: %v", err)
	}
	if got != "FINAL" {
		t.Fatalf("result = %q, want FINAL (last agent_message only)", got)
	}
}
