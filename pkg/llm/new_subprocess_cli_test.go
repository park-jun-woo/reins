//ff:func feature=llm type=helper control=sequence
//ff:what TestNewSubprocessCLI — claude/grok/codex/geminicli 이름이 각 CLI 어댑터 concrete 타입으로 디스패치되고, 미지 이름은 nil을 반환하는지 검증. 세션 env를 비워 Stateless 기본을 고정.

package llm

import "testing"

// TestNewSubprocessCLI: each known name dispatches to its concrete CLI adapter; an
// unknown name yields nil.
func TestNewSubprocessCLI(t *testing.T) {
	t.Setenv("REINS_CLAUDE_SESSION", "")
	t.Setenv("REINS_GROK_SESSION", "")
	t.Setenv("REINS_CODEX_SESSION", "")
	t.Setenv("REINS_GEMINI_SESSION", "")

	if _, ok := newSubprocessCLI("claude", "m").(*ClaudeCLI); !ok {
		t.Fatalf("claude → %T, want *ClaudeCLI", newSubprocessCLI("claude", "m"))
	}
	if _, ok := newSubprocessCLI("grok", "m").(*GrokCLI); !ok {
		t.Fatalf("grok → %T, want *GrokCLI", newSubprocessCLI("grok", "m"))
	}
	if _, ok := newSubprocessCLI("codex", "m").(*CodexCLI); !ok {
		t.Fatalf("codex → %T, want *CodexCLI", newSubprocessCLI("codex", "m"))
	}
	if _, ok := newSubprocessCLI("geminicli", "m").(*GeminiCLI); !ok {
		t.Fatalf("geminicli → %T, want *GeminiCLI", newSubprocessCLI("geminicli", "m"))
	}
	if b := newSubprocessCLI("nope", "m"); b != nil {
		t.Fatalf("unknown → %v, want nil", b)
	}
}
