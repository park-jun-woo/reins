//ff:func feature=cli type=helper control=sequence level=error
//ff:what 테스트 헬퍼. stubDef로 loop를 옵트인한 퀘스트 CLI를 session/out·args로 실행한다(newLoopRootDef 래퍼).

package cli

import (
	"testing"
)

// newLoopRoot builds a quest CLI (stubDef) with the loop opted in and runs one
// command, returning combined output.
func newLoopRoot(t *testing.T, opts Options, session, out string, args ...string) (string, error) {
	return newLoopRootDef(t, stubDef{}, opts, session, out, args...)
}
