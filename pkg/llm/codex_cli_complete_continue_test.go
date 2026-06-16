//ff:func feature=llm type=adapter control=sequence level=error
//ff:what TestCodexCLICompleteContinue — Continue 모드 운반·read-only 불변식 회귀 고정: 1차 argv는 --ephemeral 없음+`-s read-only` 있음+resume 없음, 1차 JSONL thread_id "T1"이 포인터 리시버 sid로 운반되어 2차 argv가 `exec -s read-only resume T1 …`(-s가 resume "앞"·--ephemeral 없음)이 되는지 execCodex 스텁으로 검증. -s 누락/오배치 시 테스트가 깨진다.

package llm

import (
	"context"
	"testing"
)

// TestCodexCLICompleteContinue: in Continue mode the first call carries -s read-only,
// no --ephemeral, and no resume; the response's thread_id is stored on the pointer
// receiver; the second call's argv is `exec -s read-only resume T1 …` with -s kept
// BEFORE the resume subcommand and still no --ephemeral. A dropped or misplaced -s
// breaks this test (the read-only invariant is pinned on the Continue path).
func TestCodexCLICompleteContinue(t *testing.T) {
	orig := execCodex
	defer func() { execCodex = orig }()

	var argvs [][]string
	execCodex = func(ctx context.Context, bin string, argv []string, stdin string) (string, string, error) {
		argvs = append(argvs, argv)
		return `{"type":"thread.started","thread_id":"T1"}
{"type":"item.completed","item":{"type":"agent_message","text":"OK"}}`, "", nil
	}

	c := &CodexCLI{Session: SessionMode{Kind: Continue}}

	if _, err := c.Complete("SYS", "USR1"); err != nil {
		t.Fatalf("first Complete error: %v", err)
	}
	if hasFlag(argvs[0], "--ephemeral") {
		t.Fatalf("first argv %v must not contain --ephemeral (must persist for resume)", argvs[0])
	}
	if hasFlag(argvs[0], "resume") {
		t.Fatalf("first argv %v must not contain resume (sid empty)", argvs[0])
	}
	if flagValue(argvs[0], "-s") != "read-only" {
		t.Fatalf("first argv %v missing -s read-only", argvs[0])
	}

	if _, err := c.Complete("SYS", "USR2"); err != nil {
		t.Fatalf("second Complete error: %v", err)
	}
	if hasFlag(argvs[1], "--ephemeral") {
		t.Fatalf("second argv %v must not contain --ephemeral", argvs[1])
	}
	if flagValue(argvs[1], "resume") != "T1" {
		t.Fatalf("second argv %v missing resume T1 (thread_id carry)", argvs[1])
	}
	if flagValue(argvs[1], "-s") != "read-only" {
		t.Fatalf("second argv %v missing -s read-only (read-only invariant lost on resume)", argvs[1])
	}
	if si, ri := argvIndex(argvs[1], "-s"), argvIndex(argvs[1], "resume"); si < 0 || si >= ri {
		t.Fatalf("second argv %v: -s (idx %d) must precede resume (idx %d) — exec-level placement", argvs[1], si, ri)
	}
}
