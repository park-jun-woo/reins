//ff:func feature=llm type=helper control=iteration dimension=1 level=error
//ff:what TestApplyBackendOptErrors — 미지 키와 타입 파싱 실패(int·float·bool)가 각각 에러로 거부되는지 테이블로 검증.

package llm

import "testing"

// TestApplyBackendOptErrors: an unknown key and bad typed values are rejected.
func TestApplyBackendOptErrors(t *testing.T) {
	cases := []struct{ key, val string }{
		{"nope", "1"},
		{"max_output_tokens", "abc"},
		{"num_ctx", "x"},
		{"temperature", "hot"},
		{"think", "maybe"},
	}
	for _, c := range cases {
		t.Run(c.key+"="+c.val, func(t *testing.T) {
			opts := backendOpts{present: map[string]bool{}}
			if err := applyBackendOpt(&opts, c.key, c.val); err == nil {
				t.Fatalf("applyBackendOpt(%q,%q) = nil error, want error", c.key, c.val)
			}
		})
	}
}
