//ff:func feature=llm type=helper control=sequence dimension=1 level=error
//ff:what TestWithNoToolsPreamble — 빈 system이면 withNoToolsPreamble가 프리앰블 단독을 반환하고, 비빈 system이면 프리앰블이 선행(접두)하며 system을 후행 포함함을 단언(선/후행 순서 + 결합 형태 고정).

package llm

import (
	"strings"
	"testing"
)

// TestWithNoToolsPreamble asserts both contracts of withNoToolsPreamble:
// an empty system yields the preamble alone, and a non-empty system is appended
// after the preamble (preamble leads, system follows).
func TestWithNoToolsPreamble(t *testing.T) {
	if got := withNoToolsPreamble(""); got != claudeNoToolsPreamble {
		t.Fatalf("empty system = %q, want preamble alone %q", got, claudeNoToolsPreamble)
	}

	const sys = "Emit JSON per S-13."
	got := withNoToolsPreamble(sys)
	want := claudeNoToolsPreamble + "\n\n" + sys
	if got != want {
		t.Fatalf("non-empty system = %q, want %q", got, want)
	}
	if !strings.HasPrefix(got, claudeNoToolsPreamble) {
		t.Fatalf("preamble must lead: %q", got)
	}
	if !strings.Contains(got, sys) {
		t.Fatalf("system must be contained: %q", got)
	}
}
