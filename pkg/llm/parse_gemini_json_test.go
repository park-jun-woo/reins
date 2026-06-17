//ff:func feature=llm type=adapter control=iteration dimension=1 level=error
//ff:what TestParseGeminiJSON — parseGeminiJSON을 테이블로 단언: (1) {response,stats,error:null} 봉투에서 response 추출, (2) error 봉투(non-null)는 에러, (3) json 미산출(plain 텍스트, 이슈 #11184) 폴백으로 trimmed stdout을 response로, (4) 완전히 빈 출력은 에러. wantErr 행은 err!=nil만 확인.

package llm

import "testing"

// TestParseGeminiJSON asserts parseGeminiJSON over a table: it extracts the response
// from the envelope, errors on a non-null error field, falls back to plain text when
// stdout is not the json envelope (#11184), and errors on empty output.
func TestParseGeminiJSON(t *testing.T) {
	cases := []struct {
		name     string
		in       string
		wantText string
		wantErr  bool
	}{
		{
			name:     "envelope",
			in:       `{"response":"OK","stats":{"tools":{"calls":0}},"error":null}`,
			wantText: "OK",
		},
		{
			name:    "error-envelope",
			in:      `{"response":null,"error":{"message":"boom"}}`,
			wantErr: true,
		},
		{
			name:     "fallback-plain",
			in:       "  plain text answer\n",
			wantText: "plain text answer",
		},
		{name: "empty", in: "   ", wantErr: true},
	}
	for _, tc := range cases {
		text, err := parseGeminiJSON(tc.in)
		if tc.wantErr && err == nil {
			t.Fatalf("%s: err = nil, want error", tc.name)
		}
		if tc.wantErr {
			continue
		}
		if err != nil {
			t.Fatalf("%s: unexpected err %v", tc.name, err)
		}
		if text != tc.wantText {
			t.Fatalf("%s: text = %q, want %q", tc.name, text, tc.wantText)
		}
	}
}
