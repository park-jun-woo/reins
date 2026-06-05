//ff:func feature=cli type=helper control=sequence level=error
//ff:what 테스트 헬퍼. 주어진 session/out 경로·stdin·args로 stub 퀘스트 CLI를 실행하고 합쳐진 출력을 돌려준다(명령 에러 시 테스트 실패).

package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/park-jun-woo/reins/pkg/gate"
)

// runCmd executes the stub quest CLI with the given session/out paths, stdin, and
// args, returning the combined stdout+stderr. It fails the test on a command error.
func runCmd(t *testing.T, def gate.Definition, session, out, in string, args ...string) string {
	t.Helper()
	cmd := NewQuestCmd("stub", def, Options{})
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetIn(strings.NewReader(in))
	full := append([]string{"--session", session, "--out", out}, args...)
	cmd.SetArgs(full)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute %v: %v\n%s", args, err, buf.String())
	}
	return buf.String()
}
