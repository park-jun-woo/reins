//ff:func feature=llm type=model control=sequence
//ff:what TestCallFuncComplete — 함수형 Backend가 래핑한 함수에 위임해 system/user를 전달하고 결과를 반환하는지 검증.

package llm

import (
	"testing"
)

// TestCallFuncComplete delegates to the wrapped function, passing system/user and
// returning its result.
func TestCallFuncComplete(t *testing.T) {
	var gotSys, gotUser string
	var f CallFunc = func(system, user string) (string, error) {
		gotSys, gotUser = system, user
		return "out", nil
	}
	got, err := f.Complete("S", "U")
	if err != nil {
		t.Fatalf("Complete error: %v", err)
	}
	if got != "out" {
		t.Fatalf("result = %q, want out", got)
	}
	if gotSys != "S" || gotUser != "U" {
		t.Fatalf("delegated args = (%q,%q), want (S,U)", gotSys, gotUser)
	}
}
