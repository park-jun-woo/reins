//ff:func feature=textmatch type=helper control=sequence
//ff:what Contains가 빈/공백 토큰을 false로 막아 strings.Contains(src,"") 자명-참(빈-앵커 치즈)을 차단하는지 검증한다.

package textmatch

import "testing"

func TestContainsEmptyTokenFalse(t *testing.T) {
	if Contains("anything", "   ") {
		t.Fatal("empty/whitespace token must not match (empty-anchor cheese)")
	}
}
