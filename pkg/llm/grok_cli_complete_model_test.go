//ff:func feature=llm type=adapter control=sequence level=error
//ff:what TestGrokCLICompleteModel — Model 비어있지 않으면 -m <model>을 붙이고, 빈 Model(grok:default)이면 -m을 통째로 생략(CLI 기본 모델)하는지 execGrok 스텁으로 검증.

package llm

import (
	"context"
	"testing"
)

// TestGrokCLICompleteModel: a non-empty Model adds -m <model>; an empty Model omits
// -m entirely (use the CLI default).
func TestGrokCLICompleteModel(t *testing.T) {
	orig := execGrok
	defer func() { execGrok = orig }()

	var gotArgv []string
	execGrok = func(ctx context.Context, bin string, argv []string, stdin string) (string, string, error) {
		gotArgv = argv
		return `{"text":"OK","sessionId":"S1","stopReason":"EndTurn"}`, "", nil
	}

	// Model set ⇒ -m grok-4 present.
	g := &GrokCLI{Model: "grok-4"}
	if _, err := g.Complete("SYS", "USR"); err != nil {
		t.Fatalf("Complete error: %v", err)
	}
	if v := flagValue(gotArgv, "-m"); v != "grok-4" {
		t.Fatalf("-m = %q, want grok-4", v)
	}

	// Model empty (grok:default) ⇒ -m absent.
	g2 := newGrokCLI("default")
	if _, err := g2.Complete("SYS", "USR"); err != nil {
		t.Fatalf("Complete error: %v", err)
	}
	if hasFlag(gotArgv, "-m") {
		t.Fatalf("argv %v must not contain -m when Model empty", gotArgv)
	}
}
