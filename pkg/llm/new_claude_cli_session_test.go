//ff:func feature=llm type=adapter control=iteration dimension=1
//ff:what TestNewClaudeCLISession — REINS_CLAUDE_SESSION이 세션 모드를 선택하는지 테이블로 검증: "continue"만 Continue로 옵트인하고 미설정·"stateless"·오타("contineu")는 모두 Stateless로 폴백(오타가 가짜 세션을 위조하지 못함). 본문은 assertClaudeCLISession에 위임.

package llm

import (
	"testing"
)

// TestNewClaudeCLISession: REINS_CLAUDE_SESSION selects the session mode.
// Only the exact "continue" opts in; unset, "stateless", and typos all fall
// back to Stateless so a misspelled value can never forge a session.
func TestNewClaudeCLISession(t *testing.T) {
	cases := []claudeCLISessionCase{
		{name: "unset", set: false, want: Stateless},
		{name: "stateless", set: true, val: "stateless", want: Stateless},
		{name: "typo", set: true, val: "contineu", want: Stateless},
		{name: "continue", set: true, val: "continue", want: Continue},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assertClaudeCLISession(t, tc)
		})
	}
}
