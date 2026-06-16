//ff:func feature=llm type=adapter control=sequence level=error
//ff:what TestRunSubprocessStdout — 정상 종료(printf hello)가 자식의 stdout을 반환하고 stderr는 비어있으며 err가 nil인지 실제 sh로 검증; sh 부재 시 skip.

package llm

import (
	"context"
	"os/exec"
	"testing"
)

// TestRunSubprocessStdout: a clean exit returns the child's stdout with no stderr
// and a nil error.
func TestRunSubprocessStdout(t *testing.T) {
	sh, err := exec.LookPath("sh")
	if err != nil {
		t.Skip("sh not on PATH")
	}
	stdout, stderr, err := runSubprocess(context.Background(), sh, []string{"-c", "printf hello"}, "")
	if err != nil {
		t.Fatalf("runSubprocess error: %v", err)
	}
	if stdout != "hello" {
		t.Fatalf("stdout = %q, want hello", stdout)
	}
	if stderr != "" {
		t.Fatalf("stderr = %q, want empty", stderr)
	}
}
