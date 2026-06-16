//ff:func feature=llm type=adapter control=sequence level=error
//ff:what TestCodexCLICompleteModel — Model 비어있지 않으면 `-m <model>`을 붙이고, 빈 Model(codex:default)이면 -m을 통째로 생략(CLI 기본 모델)하는지 execCodex 스텁으로 검증.

package llm

import (
	"context"
	"testing"
)

// TestCodexCLICompleteModel: a non-empty Model adds -m <model>; an empty Model
// (codex:default) omits -m entirely (use the CLI default).
func TestCodexCLICompleteModel(t *testing.T) {
	orig := execCodex
	defer func() { execCodex = orig }()

	var gotArgv []string
	execCodex = func(ctx context.Context, bin string, argv []string, stdin string) (string, string, error) {
		gotArgv = argv
		return `{"type":"item.completed","item":{"type":"agent_message","text":"OK"}}`, "", nil
	}

	// Model set ⇒ -m gpt-5 present.
	c := &CodexCLI{Model: "gpt-5"}
	if _, err := c.Complete("SYS", "USR"); err != nil {
		t.Fatalf("Complete error: %v", err)
	}
	if v := flagValue(gotArgv, "-m"); v != "gpt-5" {
		t.Fatalf("-m = %q, want gpt-5", v)
	}

	// Model empty (codex:default) ⇒ -m absent.
	c2 := newCodexCLI("default")
	if _, err := c2.Complete("SYS", "USR"); err != nil {
		t.Fatalf("Complete error: %v", err)
	}
	if hasFlag(gotArgv, "-m") {
		t.Fatalf("argv %v must not contain -m when Model empty", gotArgv)
	}
}
