//ff:func feature=cli type=helper control=sequence level=error
//ff:what 테스트 헬퍼. stubDef로 agent를 옵트인한 퀘스트 CLI를 session/out·args로 실행한다(newAgentRootDef 래퍼).

package cli

import (
	"testing"
)

// newAgentRoot builds a quest CLI (stubDef) with the agent opted in and runs one
// command, returning combined output.
func newAgentRoot(t *testing.T, opts Options, session, out string, args ...string) (string, error) {
	return newAgentRootDef(t, stubDef{}, opts, session, out, args...)
}
