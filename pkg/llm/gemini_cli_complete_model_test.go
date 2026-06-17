//ff:func feature=llm type=adapter control=sequence level=error
//ff:what TestGeminiCLICompleteModel — Model 비어있지 않으면 `-m <model>`을 붙이고, 빈 Model(geminicli:default)이면 -m을 통째로 생략(CLI 기본 모델)하는지 execGemini 스텁으로 검증.

package llm

import (
	"context"
	"testing"
)

// TestGeminiCLICompleteModel: a non-empty Model adds -m <model>; an empty Model
// (geminicli:default) omits -m entirely (use the CLI default).
func TestGeminiCLICompleteModel(t *testing.T) {
	orig := execGemini
	defer func() { execGemini = orig }()

	var gotArgv []string
	execGemini = func(ctx context.Context, bin string, argv []string, stdin string) (string, string, error) {
		gotArgv = argv
		return `{"response":"OK","error":null}`, "", nil
	}

	// Model set ⇒ -m gemini-2.5-pro present.
	c := &GeminiCLI{Model: "gemini-2.5-pro"}
	if _, err := c.Complete("SYS", "USR"); err != nil {
		t.Fatalf("Complete error: %v", err)
	}
	if v := flagValue(gotArgv, "-m"); v != "gemini-2.5-pro" {
		t.Fatalf("-m = %q, want gemini-2.5-pro", v)
	}

	// Model empty (geminicli:default) ⇒ -m absent.
	c2 := newGeminiCLI("default")
	if _, err := c2.Complete("SYS", "USR"); err != nil {
		t.Fatalf("Complete error: %v", err)
	}
	if hasFlag(gotArgv, "-m") {
		t.Fatalf("argv %v must not contain -m when Model empty", gotArgv)
	}
}
