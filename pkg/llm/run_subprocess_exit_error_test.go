//ff:func feature=llm type=adapter control=sequence level=error
//ff:what TestRunSubprocessExitError — 비정상 종료(exit 3)가 non-nil err를 표면화하면서도 자식이 stderr에 쓴 내용을 캡처하는지 실제 sh로 검증; sh 부재 시 skip.

package llm

import (
	"context"
	"os/exec"
	"strings"
	"testing"
)

// TestRunSubprocessExitError: a non-zero exit surfaces a non-nil error while still
// capturing whatever the child wrote to stderr.
func TestRunSubprocessExitError(t *testing.T) {
	sh, err := exec.LookPath("sh")
	if err != nil {
		t.Skip("sh not on PATH")
	}
	stdout, stderr, err := runSubprocess(context.Background(), sh, []string{"-c", "printf boom 1>&2; exit 3"}, "")
	if err == nil {
		t.Fatal("runSubprocess = nil error, want non-zero exit error")
	}
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	if !strings.Contains(stderr, "boom") {
		t.Fatalf("stderr = %q, want it to capture child stderr", stderr)
	}
}
