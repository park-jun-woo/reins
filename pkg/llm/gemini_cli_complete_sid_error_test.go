//ff:func feature=llm type=adapter control=sequence level=error
//ff:what TestGeminiCLICompleteSIDError — Continue 1차에서 UUID 발급(geminiRandRead seam)이 실패하면 Complete가 서브프로세스를 돌리기 전에 그 에러를 표면화하는지 검증(execGemini는 호출되지 않아야 함).

package llm

import (
	"context"
	"errors"
	"testing"
)

// TestGeminiCLICompleteSIDError: in Continue mode the first call issues a UUID; if the
// randomness seam fails, Complete surfaces that error before running the subprocess.
func TestGeminiCLICompleteSIDError(t *testing.T) {
	origRand := geminiRandRead
	defer func() { geminiRandRead = origRand }()
	geminiRandRead = func(b []byte) (int, error) { return 0, errors.New("no entropy") }

	origExec := execGemini
	defer func() { execGemini = origExec }()
	called := false
	execGemini = func(ctx context.Context, bin string, argv []string, stdin string) (string, string, error) {
		called = true
		return `{"response":"OK"}`, "", nil
	}

	c := &GeminiCLI{Session: SessionMode{Kind: Continue}}
	if _, err := c.Complete("SYS", "USR"); err == nil {
		t.Fatal("Complete = nil error, want UUID-issue error")
	}
	if called {
		t.Fatal("execGemini was called; Complete must fail before spawning the subprocess")
	}
}
