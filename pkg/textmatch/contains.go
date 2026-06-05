//ff:func feature=textmatch type=helper control=sequence
//ff:what token이 source의 substring인지 판정(둘 다 Normalize 후). 정규화 후 빈/공백 토큰은 false — strings.Contains(src,"")의 자명-참 함정(빈-앵커 치즈 벡터) 차단. 순수 함수.

package textmatch

import "strings"

// Contains reports whether token appears as a substring of source after both are
// Normalize'd. An empty/whitespace-only token returns false — never the trivially
// true strings.Contains(src, "") (the empty-anchor cheese vector; reins Phase005,
// distilled from ccnews Phase009 L0).
func Contains(source, token string) bool {
	nt := Normalize(token)
	if nt == "" {
		return false
	}
	return strings.Contains(Normalize(source), nt)
}
