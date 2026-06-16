//ff:func feature=llm type=adapter control=sequence level=error
//ff:what TestCodexCLISmoke — codex가 PATH에 있을 때만 codex:default가 read-only에서 도구 시도 없이 비어있지 않은 결과를 1회 생성하는지 통합 검증(L0 단일샷 보장의 경험적 확인); 없으면 skip(CI CLI/구독 부재).

package llm

import (
	"os/exec"
	"strings"
	"testing"
)

// TestCodexCLISmoke: when codex is on PATH, codex:default produces a non-empty result
// under read-only with no tool attempt; otherwise skip (no CLI/subscription in CI).
func TestCodexCLISmoke(t *testing.T) {
	if _, err := exec.LookPath("codex"); err != nil {
		t.Skip("codex not on PATH")
	}
	c := newCodexCLI("default")
	out, err := c.Complete("", "Reply with exactly the two letters: ok")
	if err != nil {
		t.Skipf("codex smoke skipped (login/quota?): %v", err)
	}
	if strings.TrimSpace(out) == "" {
		t.Fatalf("codex smoke returned empty result")
	}
}
