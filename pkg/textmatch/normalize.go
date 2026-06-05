//ff:func feature=textmatch type=helper control=sequence
//ff:what 표면형 정규화. 유니코드 NFC 적용 후 공백 런을 단일 스페이스로 접고 양끝을 트림한다. 원천·토큰 양쪽에 동일 적용(매핑 추론 아님). NFC가 합성/분해형 분음부호를 통일해 false-negative를 막는다. 순수 함수.

package textmatch

import (
	"strings"

	"golang.org/x/text/unicode/norm"
)

// Normalize applies Unicode NFC, collapses every run of whitespace to a single
// space, and trims the ends. It is applied identically to the source text and to
// each token so matching compares like-for-like surface forms — never an inferred
// mapping. NFC unifies composed/decomposed forms (é = "e"+◌́ vs "é"), which is what
// closes the diacritic/combining-mark false-negatives a whitespace-only normalize
// leaves open.
func Normalize(s string) string {
	return strings.Join(strings.Fields(norm.NFC.String(s)), " ")
}
