//ff:func feature=textmatch type=helper control=sequence
//ff:what Contains가 NFC 정규화로 합성/분해형 토큰(café)을 동일 substring으로 일치시키는지 검증한다.

package textmatch

import "testing"

func TestContainsNFC(t *testing.T) {
	// "é" decomposed (e + combining acute) vs composed should match after NFC.
	source := "café society"
	if !Contains(source, "café") {
		t.Fatal("NFC normalization should make composed/decomposed forms equal")
	}
}
