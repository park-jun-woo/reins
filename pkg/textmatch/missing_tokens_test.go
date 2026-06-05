//ff:func feature=textmatch type=helper control=sequence
//ff:what MissingTokens가 원천에 없는 토큰과 빈/공백 토큰을 missing으로 모으고 원래 표면형 순서를 보존하는지 검증한다.

package textmatch

import "testing"

func TestMissingTokens(t *testing.T) {
	miss := MissingTokens("the quick brown fox", []string{"quick", "lazy", ""})
	if len(miss) != 2 || miss[0] != "lazy" || miss[1] != "" {
		t.Fatalf("missing = %v", miss)
	}
}
