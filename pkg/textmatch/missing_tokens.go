//ff:func feature=textmatch type=helper control=iteration dimension=1
//ff:what tokens 중 source에 없는 것들을 반환(다건 앵커 검사). source는 한 번만 정규화하고, 정규화 후 빈/공백 토큰은 missing으로 센다. 규칙이 첫 위반 토큰으로 Fact 만들기 좋게 원래 표면형을 돌려준다. 순수 함수.

package textmatch

import "strings"

// MissingTokens returns the tokens that are not present in source (after
// normalization). Source is normalized once. Empty/whitespace tokens count as
// missing. Handy for a gate rule to build a Fact naming the first offender.
func MissingTokens(source string, tokens []string) []string {
	ns := Normalize(source)
	var miss []string
	for _, t := range tokens {
		nt := Normalize(t)
		if nt == "" || !strings.Contains(ns, nt) {
			miss = append(miss, t)
		}
	}
	return miss
}
