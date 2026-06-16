//ff:func feature=llm type=adapter control=sequence level=error
//ff:what TestGrokCLICompleteContinue — Continue 모드에서 1차 Complete argv에 -r가 없고(빈 sid), 응답 sessionId:"S1"이 포인터 리시버 sid에 운반되어 2차 argv에 -r S1이 실리는지 execGrok 스텁으로 검증.

package llm

import (
	"context"
	"testing"
)

// TestGrokCLICompleteContinue: in Continue mode the first call carries no -r (empty
// sid), the response's sessionId is stored on the pointer receiver, and the second
// call's argv carries -r S1.
func TestGrokCLICompleteContinue(t *testing.T) {
	orig := execGrok
	defer func() { execGrok = orig }()

	var argvs [][]string
	execGrok = func(ctx context.Context, bin string, argv []string, stdin string) (string, string, error) {
		argvs = append(argvs, argv)
		return `{"text":"OK","sessionId":"S1","stopReason":"EndTurn"}`, "", nil
	}

	g := &GrokCLI{Session: SessionMode{Kind: Continue}}

	if _, err := g.Complete("SYS", "USR1"); err != nil {
		t.Fatalf("first Complete error: %v", err)
	}
	if hasFlag(argvs[0], "-r") {
		t.Fatalf("first argv %v must not contain -r (sid empty)", argvs[0])
	}

	if _, err := g.Complete("SYS", "USR2"); err != nil {
		t.Fatalf("second Complete error: %v", err)
	}
	if v := flagValue(argvs[1], "-r"); v != "S1" {
		t.Fatalf("-r = %q, want S1", v)
	}
}
