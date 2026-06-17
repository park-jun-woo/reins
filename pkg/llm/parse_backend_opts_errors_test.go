//ff:func feature=llm type=helper control=iteration dimension=1 level=error
//ff:what TestParseBackendOptsErrors — '=' 없는 세그먼트·미지 키·int/float/bool 파싱 실패가 각각 에러로 거부되는지(loud, no silent) 테이블로 검증.

package llm

import "testing"

// TestParseBackendOptsErrors: a missing '=', an unknown key, and bad int/float/bool
// values are each rejected with an error.
func TestParseBackendOptsErrors(t *testing.T) {
	cases := []string{
		"foo",                   // no '='
		"bar=1",                 // unknown key
		"max_output_tokens=abc", // bad int
		"num_ctx=x",             // bad int
		"temperature=hot",       // bad float
		"think=maybe",           // bad bool
	}
	for _, raw := range cases {
		t.Run(raw, func(t *testing.T) {
			if _, err := parseBackendOpts(raw); err == nil {
				t.Fatalf("parseBackendOpts(%q) = nil error, want error", raw)
			}
		})
	}
}
