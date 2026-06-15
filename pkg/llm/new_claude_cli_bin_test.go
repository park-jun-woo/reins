//ff:func feature=llm type=adapter control=sequence
//ff:what TestNewClaudeCLIBin — REINS_CLAUDE_BIN env가 newClaudeCLI의 Bin을 덮어쓰는지 t.Setenv로 검증. 무서브프로세스.

package llm

import (
	"testing"
)

// TestNewClaudeCLIBin: REINS_CLAUDE_BIN overrides Bin.
func TestNewClaudeCLIBin(t *testing.T) {
	t.Setenv("REINS_CLAUDE_BIN", "/opt/claude")
	if c := newClaudeCLI("opus"); c.Bin != "/opt/claude" {
		t.Fatalf("Bin = %q, want /opt/claude", c.Bin)
	}
}
