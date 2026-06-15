//ff:func feature=llm type=adapter control=sequence
//ff:what assertClaudeCLISession — 한 세션 케이스를 검증하는 헬퍼: REINS_CLAUDE_SESSION을 설정/공백처리한 뒤 newClaudeCLI의 Session.Kind가 기대값과 같은지 단언(TestNewClaudeCLISession의 루프 본문 추출, Q4 회피).

package llm

import (
	"testing"
)

// assertClaudeCLISession applies one case: it sets (or blanks) the env var and
// checks that newClaudeCLI selects the expected session kind. Blanking on the
// unset case ensures a leaked outer env cannot skew the test.
func assertClaudeCLISession(t *testing.T, tc claudeCLISessionCase) {
	t.Helper()
	if tc.set {
		t.Setenv("REINS_CLAUDE_SESSION", tc.val)
	} else {
		t.Setenv("REINS_CLAUDE_SESSION", "")
	}
	c := newClaudeCLI("opus")
	if c.Session.Kind != tc.want {
		t.Fatalf("Session.Kind = %v, want %v", c.Session.Kind, tc.want)
	}
}
