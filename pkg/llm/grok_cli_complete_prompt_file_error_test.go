//ff:func feature=llm type=adapter control=sequence level=error
//ff:what TestGrokCLICompletePromptFileError — 프롬프트 임시파일을 만들 수 없을 때(TMPDIR가 없는 디렉터리) Complete가 execGrok에 닿기 전 "grok prompt file" 에러로 조기 실패하고 execGrok이 호출되지 않는지 검증.

package llm

import (
	"context"
	"path/filepath"
	"strings"
	"testing"
)

// TestGrokCLICompletePromptFileError: when the prompt temp file cannot be created
// (TMPDIR points at a nonexistent directory) Complete fails early — before execGrok
// is ever reached — with a "grok prompt file" error.
func TestGrokCLICompletePromptFileError(t *testing.T) {
	orig := execGrok
	defer func() { execGrok = orig }()
	called := false
	execGrok = func(ctx context.Context, bin string, argv []string, stdin string) (string, string, error) {
		called = true
		return `{"text":"OK","stopReason":"EndTurn"}`, "", nil
	}

	// os.CreateTemp("", …) resolves the dir via os.TempDir(), which honors $TMPDIR
	// on unix; pointing it at a missing path forces the create to fail.
	t.Setenv("TMPDIR", filepath.Join(t.TempDir(), "no-such-dir"))

	g := &GrokCLI{}
	_, err := g.Complete("SYS", "USR")
	if err == nil {
		t.Fatal("Complete = nil error, want prompt-file create error")
	}
	if !strings.Contains(err.Error(), "grok prompt file") {
		t.Fatalf("error = %q, want it to mention the prompt file", err.Error())
	}
	if called {
		t.Fatal("execGrok must not run when the prompt file cannot be created")
	}
}
