//ff:func feature=cli type=helper control=sequence level=error
//ff:what 테스트 헬퍼. 명시 Definition으로 loop를 옵트인한 stub 퀘스트 CLI를 만들어 session/out·args로 1회 실행하고 합쳐진 출력을 돌려준다.

package cli

import (
	"bytes"
	"testing"

	"github.com/park-jun-woo/reins/pkg/gate"
)

// newLoopRootDef builds a quest CLI (explicit Definition) with the loop opted in
// and runs one command, returning combined output.
func newLoopRootDef(t *testing.T, def gate.Definition, opts Options, session, out string, args ...string) (string, error) {
	t.Helper()
	cmd := NewQuestCmd("stub", def, opts)
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	full := append([]string{"--session", session, "--out", out}, args...)
	cmd.SetArgs(full)
	err := cmd.Execute()
	return buf.String(), err
}
