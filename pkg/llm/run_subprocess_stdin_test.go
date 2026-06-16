//ff:func feature=llm type=adapter control=sequence level=error
//ff:what TestRunSubprocessStdin — stdin이 자식에게 먹여져 stdout으로 되돌아오는지(cat) 실제 바이너리로 검증해 stdin 운반선이 연결됐음을 증명; cat 부재 시 skip.

package llm

import (
	"context"
	"os/exec"
	"testing"
)

// TestRunSubprocessStdin: stdin is fed to the child and echoed back on stdout,
// proving the stdin carrier is wired up.
func TestRunSubprocessStdin(t *testing.T) {
	cat, err := exec.LookPath("cat")
	if err != nil {
		t.Skip("cat not on PATH")
	}
	stdout, _, err := runSubprocess(context.Background(), cat, nil, "piped-in")
	if err != nil {
		t.Fatalf("runSubprocess error: %v", err)
	}
	if stdout != "piped-in" {
		t.Fatalf("stdout = %q, want piped-in (stdin must reach the child)", stdout)
	}
}
