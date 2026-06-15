//ff:func feature=llm type=adapter control=sequence level=error
//ff:what TestClaudeCLICompleteContinue — Continue 모드에서 1차 Complete argv에 --resume가 없고(빈 sid), 응답 session_id:"S1"이 포인터 리시버 sid에 운반되어 2차 argv에 --resume S1이 실리는지 execClaude 스텁으로 검증.

package llm

import (
	"context"
	"testing"
)

// TestClaudeCLICompleteContinue: in Continue mode the first call carries no
// --resume (empty sid), the response's session_id is stored on the pointer
// receiver, and the second call's argv carries --resume S1.
func TestClaudeCLICompleteContinue(t *testing.T) {
	orig := execClaude
	defer func() { execClaude = orig }()

	var argvs [][]string
	execClaude = func(ctx context.Context, bin string, argv []string, stdin string) (string, string, error) {
		argvs = append(argvs, argv)
		return `{"result":"OK","session_id":"S1","is_error":false}`, "", nil
	}

	c := &ClaudeCLI{Session: SessionMode{Kind: Continue}}

	if _, err := c.Complete("SYS", "USR1"); err != nil {
		t.Fatalf("first Complete error: %v", err)
	}
	if hasFlag(argvs[0], "--resume") {
		t.Fatalf("first argv %v must not contain --resume (sid empty)", argvs[0])
	}

	if _, err := c.Complete("SYS", "USR2"); err != nil {
		t.Fatalf("second Complete error: %v", err)
	}
	if !hasFlag(argvs[1], "--resume") {
		t.Fatalf("second argv %v missing --resume", argvs[1])
	}
	if v := flagValue(argvs[1], "--resume"); v != "S1" {
		t.Fatalf("--resume = %q, want S1", v)
	}
}
