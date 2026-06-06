//ff:func feature=llm type=model control=sequence level=error
//ff:what TestCallFuncCompleteError — 래핑한 함수가 반환한 에러를 그대로 전파하는지 검증.

package llm

import (
	"errors"
	"testing"
)

// TestCallFuncCompleteError propagates the wrapped function's error.
func TestCallFuncCompleteError(t *testing.T) {
	want := errors.New("boom")
	var f CallFunc = func(string, string) (string, error) { return "", want }
	if _, err := f.Complete("a", "b"); !errors.Is(err, want) {
		t.Fatalf("error = %v, want %v", err, want)
	}
}
