//ff:func feature=cli type=helper control=sequence
//ff:what TestComposeSystemFallbackUsed — 빈 전역·빈 코칭일 때 fallback이 reins 기본 가이던스("deterministic gate")를 담는지 검증.

package cli

import (
	"strings"
	"testing"
)

// TestComposeSystemFallbackUsed: the fallback contains the canonical reins guidance.
func TestComposeSystemFallbackUsed(t *testing.T) {
	got := composeSystem("", "")
	if !strings.Contains(got, "deterministic gate") {
		t.Fatalf("fallback = %q, want generic reins prompt", got)
	}
}
