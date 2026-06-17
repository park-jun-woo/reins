//ff:func feature=llm type=adapter control=sequence level=error
//ff:what TestGeminiCLICompleteContinue — Continue 모드 세션 운반 회귀 고정: 1차 argv에 reins 발급 UUID로 `--session-id <uuid>`(v4 포맷)가 있고 `--resume` 없음, 그 UUID가 포인터 리시버 sid로 운반되어 2차 argv가 `--resume latest`(--session-id 없음)가 되는지 execGemini 스텁으로 검증. resume/session-id 누락 시 테스트가 깨진다.

package llm

import (
	"context"
	"testing"
)

// TestGeminiCLICompleteContinue: in Continue mode the first call carries
// --session-id <reins-issued v4 uuid> and no --resume; the uuid is stored on the
// pointer receiver; the second call's argv is `--resume latest` with no --session-id.
func TestGeminiCLICompleteContinue(t *testing.T) {
	orig := execGemini
	defer func() { execGemini = orig }()

	var argvs [][]string
	execGemini = func(ctx context.Context, bin string, argv []string, stdin string) (string, string, error) {
		argvs = append(argvs, argv)
		return `{"response":"OK","error":null}`, "", nil
	}

	c := &GeminiCLI{Session: SessionMode{Kind: Continue}}

	if _, err := c.Complete("SYS", "USR1"); err != nil {
		t.Fatalf("first Complete error: %v", err)
	}
	first := flagValue(argvs[0], "--session-id")
	if !uuidV4Re.MatchString(first) {
		t.Fatalf("first argv %v: --session-id %q is not a v4 uuid", argvs[0], first)
	}
	if hasFlag(argvs[0], "--resume") {
		t.Fatalf("first argv %v must not contain --resume (session not yet created)", argvs[0])
	}
	if c.sid != first {
		t.Fatalf("sid = %q, want carried uuid %q", c.sid, first)
	}

	if _, err := c.Complete("SYS", "USR2"); err != nil {
		t.Fatalf("second Complete error: %v", err)
	}
	if flagValue(argvs[1], "--resume") != "latest" {
		t.Fatalf("second argv %v missing --resume latest", argvs[1])
	}
	if hasFlag(argvs[1], "--session-id") {
		t.Fatalf("second argv %v must not contain --session-id (resume instead)", argvs[1])
	}
}
