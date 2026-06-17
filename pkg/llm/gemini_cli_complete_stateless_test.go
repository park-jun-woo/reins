//ff:func feature=llm type=adapter control=sequence level=error
//ff:what TestGeminiCLICompleteStateless — Stateless 기본에서 gemini argv에 `-p ""`·`-o json`·`--approval-mode plan`이 있고 `--resume`·`--session-id`가 **없으며**, stdin이 `프리앰블+system+user` 선결합이고, json 봉투에서 response를 반환하는지 execGemini 스텁으로 검증.

package llm

import (
	"context"
	"testing"
)

// TestGeminiCLICompleteStateless: the default Stateless mode builds a gemini argv with
// -p "", -o json, and --approval-mode plan but no --resume / --session-id, sends the
// no-tools preamble + system + user on stdin, and returns the response from the json
// envelope.
func TestGeminiCLICompleteStateless(t *testing.T) {
	orig := execGemini
	defer func() { execGemini = orig }()

	var gotArgv []string
	var gotStdin string
	execGemini = func(ctx context.Context, bin string, argv []string, stdin string) (string, string, error) {
		gotArgv, gotStdin = argv, stdin
		return `{"response":"OK","error":null}`, "", nil
	}

	c := &GeminiCLI{}
	got, err := c.Complete("SYS", "USR")
	if err != nil {
		t.Fatalf("Complete error: %v", err)
	}
	if got != "OK" {
		t.Fatalf("result = %q, want OK", got)
	}
	if want := withNoToolsPreamble("SYS") + "\n\n" + "USR"; gotStdin != want {
		t.Fatalf("stdin = %q, want %q", gotStdin, want)
	}
	if flagValue(gotArgv, "-o") != "json" {
		t.Fatalf("argv %v missing -o json", gotArgv)
	}
	if flagValue(gotArgv, "--approval-mode") != "plan" {
		t.Fatalf("argv %v missing --approval-mode plan", gotArgv)
	}
	if !hasFlag(gotArgv, "-p") {
		t.Fatalf("argv %v missing -p", gotArgv)
	}
	if hasFlag(gotArgv, "--resume") || hasFlag(gotArgv, "--session-id") {
		t.Fatalf("argv %v must not contain --resume/--session-id in Stateless", gotArgv)
	}
	if hasFlag(gotArgv, "--all-files") || hasFlag(gotArgv, "yolo") {
		t.Fatalf("argv %v must never contain --all-files/yolo", gotArgv)
	}
}
