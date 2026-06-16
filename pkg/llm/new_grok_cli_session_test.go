//ff:func feature=llm type=adapter control=iteration dimension=1
//ff:what TestNewGrokCLISessionFallback — REINS_GROK_SESSION이 정확히 "continue"가 아닌 값(미설정·"stateless"·오타)은 모두 Stateless로 폴백하는지 테이블로 검증(오타가 가짜 세션을 위조하지 못함). 본문은 assertGrokCLISessionFallback에 위임.

package llm

import (
	"testing"
)

// TestNewGrokCLISessionFallback: any REINS_GROK_SESSION value other than the exact
// "continue" (unset, "stateless", a typo) falls back to Stateless, so a misspelled
// value can never silently forge a session.
func TestNewGrokCLISessionFallback(t *testing.T) {
	cases := []grokCLISessionFallbackCase{
		{name: "unset", set: false},
		{name: "stateless", set: true, val: "stateless"},
		{name: "typo", set: true, val: "contineu"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assertGrokCLISessionFallback(t, tc)
		})
	}
}
