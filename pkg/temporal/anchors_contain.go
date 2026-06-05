//ff:func feature=temporal type=helper control=iteration dimension=1
//ff:what 토큰이 anchors 중 하나라도 substring으로 들어있는지 검사한다(textmatch.Contains 재사용, 정규화 위임). 하나라도 포함하면 true.

package temporal

import "github.com/park-jun-woo/reins/pkg/textmatch"

// anchorsContain reports whether token appears in any of the anchors.
func anchorsContain(anchors []string, token string) bool {
	for _, a := range anchors {
		if textmatch.Contains(a, token) {
			return true
		}
	}
	return false
}
