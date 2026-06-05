//ff:func feature=cli type=helper control=sequence level=error
//ff:what 테스트 헬퍼. runCmd처럼 실행하되 명령 에러를 기대한다(성공하면 테스트 실패, 아니면 에러 반환). 서브커맨드의 level=error 분기 자극용.

package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/park-jun-woo/reins/pkg/gate"
)

// runCmdErr executes the stub quest CLI like runCmd but expects a command error; it
// fails the test if the command succeeds, and returns the error otherwise. Used to
// exercise the level=error branches of the subcommands.
func runCmdErr(t *testing.T, def gate.Definition, session, out, in string, args ...string) error {
	t.Helper()
	cmd := NewQuestCmd("stub", def, Options{})
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetIn(strings.NewReader(in))
	full := append([]string{"--session", session, "--out", out}, args...)
	cmd.SetArgs(full)
	err := cmd.Execute()
	if err == nil {
		t.Fatalf("execute %v: want error, got success\n%s", args, buf.String())
	}
	return err
}
