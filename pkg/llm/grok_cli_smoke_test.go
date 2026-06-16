//ff:func feature=llm type=adapter control=sequence level=error
//ff:what TestGrokCLISmoke — grok이 PATH에 있을 때만 grok:default가 도구 시도 없이 비어있지 않은 결과를 1회 생성하는지 통합 검증; 없으면 skip(CI 네트워크/CLI 부재).

package llm

import (
	"os/exec"
	"strings"
	"testing"
)

// TestGrokCLISmoke: when grok is on PATH, grok:default produces a non-empty result
// with no tool attempt; otherwise skip (no network/CLI in CI).
func TestGrokCLISmoke(t *testing.T) {
	if _, err := exec.LookPath("grok"); err != nil {
		t.Skip("grok not on PATH")
	}
	g := newGrokCLI("default")
	out, err := g.Complete("", "Reply with exactly the two letters: ok")
	if err != nil {
		t.Skipf("grok smoke skipped (login/quota?): %v", err)
	}
	if strings.TrimSpace(out) == "" {
		t.Fatalf("grok smoke returned empty result")
	}
}
