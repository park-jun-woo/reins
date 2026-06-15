//ff:func feature=llm type=adapter control=iteration dimension=1 level=error
//ff:what TestClaudeCLICompleteStateless — Stateless 기본에서 argv에 --no-session-persistence·--output-format json·--append-system-prompt·--max-turns 1·--permission-mode dontAsk가 있고 --bare·--resume는 없으며 stdin이 user 프롬프트임을 execClaude 스텁으로 검증(무서브프로세스). 정상 JSON이 result를 반환.

package llm

import (
	"context"
	"testing"
)

// TestClaudeCLICompleteStateless: the default Stateless mode builds the plain
// `claude -p` argv (--no-session-persistence, --output-format json,
// --append-system-prompt, --max-turns 1, --permission-mode dontAsk) without
// --bare or --resume, feeds the user prompt on stdin, and returns result.
func TestClaudeCLICompleteStateless(t *testing.T) {
	orig := execClaude
	defer func() { execClaude = orig }()

	var gotArgv []string
	var gotStdin string
	execClaude = func(ctx context.Context, bin string, argv []string, stdin string) (string, string, error) {
		gotArgv = argv
		gotStdin = stdin
		return `{"result":"OK","session_id":"S1","is_error":false}`, "", nil
	}

	c := &ClaudeCLI{}
	got, err := c.Complete("SYS", "USR")
	if err != nil {
		t.Fatalf("Complete error: %v", err)
	}
	if got != "OK" {
		t.Fatalf("result = %q, want OK", got)
	}
	if gotStdin != "USR" {
		t.Fatalf("stdin = %q, want USR", gotStdin)
	}
	for _, want := range []string{"--no-session-persistence", "--output-format", "--append-system-prompt", "--max-turns", "--permission-mode"} {
		if !hasFlag(gotArgv, want) {
			t.Fatalf("argv %v missing %q", gotArgv, want)
		}
	}
	if v := flagValue(gotArgv, "--output-format"); v != "json" {
		t.Fatalf("--output-format = %q, want json", v)
	}
	if v := flagValue(gotArgv, "--append-system-prompt"); v != "SYS" {
		t.Fatalf("--append-system-prompt = %q, want SYS", v)
	}
	if v := flagValue(gotArgv, "--max-turns"); v != "1" {
		t.Fatalf("--max-turns = %q, want 1", v)
	}
	if v := flagValue(gotArgv, "--permission-mode"); v != "dontAsk" {
		t.Fatalf("--permission-mode = %q, want dontAsk", v)
	}
	if hasFlag(gotArgv, "--bare") {
		t.Fatalf("argv %v must not contain --bare", gotArgv)
	}
	if hasFlag(gotArgv, "--resume") {
		t.Fatalf("argv %v must not contain --resume in Stateless", gotArgv)
	}
}
