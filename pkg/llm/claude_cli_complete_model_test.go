//ff:func feature=llm type=adapter control=sequence level=error
//ff:what TestClaudeCLICompleteModel — Model 설정 시 argv에 --model <model>이 실리고, Model 빈 값이면 --model이 부재함을 execClaude 스텁으로 검증.

package llm

import (
	"context"
	"testing"
)

// TestClaudeCLICompleteModel: a non-empty Model adds --model <model>; an empty
// Model omits --model entirely (use the CLI default).
func TestClaudeCLICompleteModel(t *testing.T) {
	orig := execClaude
	defer func() { execClaude = orig }()

	var gotArgv []string
	execClaude = func(ctx context.Context, bin string, argv []string, stdin string) (string, string, error) {
		gotArgv = argv
		return `{"result":"OK","session_id":"S1","is_error":false}`, "", nil
	}

	// Model set ⇒ --model opus present.
	c := &ClaudeCLI{Model: "opus"}
	if _, err := c.Complete("SYS", "USR"); err != nil {
		t.Fatalf("Complete error: %v", err)
	}
	if !hasFlag(gotArgv, "--model") {
		t.Fatalf("argv %v missing --model", gotArgv)
	}
	if v := flagValue(gotArgv, "--model"); v != "opus" {
		t.Fatalf("--model = %q, want opus", v)
	}

	// Model empty ⇒ --model absent.
	c2 := &ClaudeCLI{}
	if _, err := c2.Complete("SYS", "USR"); err != nil {
		t.Fatalf("Complete error: %v", err)
	}
	if hasFlag(gotArgv, "--model") {
		t.Fatalf("argv %v must not contain --model when Model empty", gotArgv)
	}
}
