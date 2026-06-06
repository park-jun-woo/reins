//ff:func feature=llm type=model control=sequence
//ff:what TestCallFuncSatisfiesBackend — CallFunc가 Backend 인터페이스로 사용 가능한지 검증.

package llm

import (
	"testing"
)

// TestCallFuncSatisfiesBackend: CallFunc is usable as a Backend.
func TestCallFuncSatisfiesBackend(t *testing.T) {
	var b Backend = CallFunc(func(string, string) (string, error) { return "ok", nil })
	if got, _ := b.Complete("", ""); got != "ok" {
		t.Fatalf("Backend.Complete = %q, want ok", got)
	}
}
