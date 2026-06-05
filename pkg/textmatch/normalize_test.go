//ff:func feature=textmatch type=helper control=iteration dimension=1
//ff:what Normalize가 공백 런 접기·양끝 트림·유니코드 공백 처리, 그리고 NFC 합성/분해형 통일(café 분해형→합성형)을 케이스별로 검증한다.

package textmatch

import "testing"

func TestNormalize(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"empty", "", ""},
		{"only spaces", "   ", ""},
		{"trim ends", "  hello  ", "hello"},
		{"collapse runs", "a   b\t\tc", "a b c"},
		{"newlines and tabs", "a\n\tb\r\nc", "a b c"},
		{"unicode space", "a b　c", "a b c"},
		{"already normal", "a b c", "a b c"},
		// "café" decomposed (e + combining acute U+0301) must become the composed
		// NFC form so it equals a composed "café" literal.
		{"nfc compose diacritic", "café", "café"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := Normalize(tc.in); got != tc.want {
				t.Errorf("Normalize(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}
