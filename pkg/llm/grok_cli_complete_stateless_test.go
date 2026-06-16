//ff:func feature=llm type=adapter control=iteration dimension=1 level=error
//ff:what TestGrokCLICompleteStateless — Stateless 기본에서 grok argv(--prompt-file<임시파일>+그 파일에 user·--output-format json·--max-turns 1·--permission-mode dontAsk·--rules 프리앰블·--tools "")를 짓고 -r 없이 stdin 빈값으로 보내며 text를 반환하는지 execGrok 스텁으로 검증.

package llm

import (
	"context"
	"os"
	"testing"
)

// TestGrokCLICompleteStateless: the default Stateless mode builds the grok argv
// (--prompt-file <temp>, --output-format json, --max-turns 1, --permission-mode
// dontAsk, --rules <preamble+SYS>, --tools "") without -r, sends nothing on stdin,
// carries the user prompt in the prompt-file, and returns text.
func TestGrokCLICompleteStateless(t *testing.T) {
	orig := execGrok
	defer func() { execGrok = orig }()

	var gotArgv []string
	var gotStdin string
	var promptFileBody string
	execGrok = func(ctx context.Context, bin string, argv []string, stdin string) (string, string, error) {
		gotArgv = argv
		gotStdin = stdin
		if pf := flagValue(argv, "--prompt-file"); pf != "" {
			if b, err := os.ReadFile(pf); err == nil {
				promptFileBody = string(b)
			}
		}
		return `{"text":"OK","sessionId":"S1","stopReason":"EndTurn"}`, "", nil
	}

	g := &GrokCLI{}
	got, err := g.Complete("SYS", "USR")
	if err != nil {
		t.Fatalf("Complete error: %v", err)
	}
	if got != "OK" {
		t.Fatalf("result = %q, want OK", got)
	}
	if gotStdin != "" {
		t.Fatalf("stdin = %q, want empty (user travels via --prompt-file)", gotStdin)
	}
	if promptFileBody != "USR" {
		t.Fatalf("prompt-file body = %q, want USR", promptFileBody)
	}
	for _, want := range []string{"--prompt-file", "--output-format", "--max-turns", "--permission-mode", "--rules", "--tools", "--disable-web-search", "--no-subagents", "--no-memory"} {
		if !hasFlag(gotArgv, want) {
			t.Fatalf("argv %v missing %q", gotArgv, want)
		}
	}
	if v := flagValue(gotArgv, "--output-format"); v != "json" {
		t.Fatalf("--output-format = %q, want json", v)
	}
	if v := flagValue(gotArgv, "--max-turns"); v != "1" {
		t.Fatalf("--max-turns = %q, want 1", v)
	}
	if v := flagValue(gotArgv, "--permission-mode"); v != "dontAsk" {
		t.Fatalf("--permission-mode = %q, want dontAsk", v)
	}
	if v := flagValue(gotArgv, "--rules"); v != withNoToolsPreamble("SYS") {
		t.Fatalf("--rules = %q, want %q", v, withNoToolsPreamble("SYS"))
	}
	if v := flagValue(gotArgv, "--tools"); v != "" {
		t.Fatalf("--tools = %q, want empty allow-list", v)
	}
	if hasFlag(gotArgv, "-r") {
		t.Fatalf("argv %v must not contain -r in Stateless", gotArgv)
	}
}
