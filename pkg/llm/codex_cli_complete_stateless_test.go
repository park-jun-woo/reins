//ff:func feature=llm type=adapter control=iteration dimension=1 level=error
//ff:what TestCodexCLICompleteStateless — Stateless 기본에서 codex argv(`exec --ephemeral -s read-only --json --skip-git-repo-check --ignore-user-config --ignore-rules -`)를 짓고, stdin이 `프리앰블+system+user` 선결합이며, JSONL에서 agent_message 텍스트를 반환하는지 execCodex 스텁으로 검증. resume/--ephemeral 위치도 단언.

package llm

import (
	"context"
	"testing"
)

// TestCodexCLICompleteStateless: the default Stateless mode builds the codex argv
// (exec --ephemeral -s read-only --json --skip-git-repo-check --ignore-user-config
// --ignore-rules -), sends the no-tools preamble + system + user on stdin, and
// returns the agent_message text parsed from the JSONL stream.
func TestCodexCLICompleteStateless(t *testing.T) {
	orig := execCodex
	defer func() { execCodex = orig }()

	var gotArgv []string
	var gotStdin string
	execCodex = func(ctx context.Context, bin string, argv []string, stdin string) (string, string, error) {
		gotArgv, gotStdin = argv, stdin
		return `{"type":"thread.started","thread_id":"T1"}
{"type":"item.completed","item":{"type":"agent_message","text":"OK"}}
{"type":"turn.completed"}`, "", nil
	}

	c := &CodexCLI{}
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
	want := []string{"exec", "--ephemeral", "-s", "read-only", "--json", "--skip-git-repo-check", "--ignore-user-config", "--ignore-rules", "-"}
	if len(gotArgv) != len(want) {
		t.Fatalf("argv = %v, want %v", gotArgv, want)
	}
	for i := range want {
		if gotArgv[i] != want[i] {
			t.Fatalf("argv[%d] = %q, want %q (full %v)", i, gotArgv[i], want[i], gotArgv)
		}
	}
}
