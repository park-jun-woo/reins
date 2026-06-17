//ff:func feature=llm type=helper control=iteration dimension=1
//ff:what TestSplitModelQuery — '?' 없음(raw 빈 문자열), '?' 있음(model|raw 분리), 콜론 포함 model 보존(gemma4:e4b?…)을 테이블로 검증.

package llm

import "testing"

// TestSplitModelQuery checks the first-'?' split, the no-query case, and colon
// preservation in the model name.
func TestSplitModelQuery(t *testing.T) {
	cases := []struct {
		in        string
		wantModel string
		wantRaw   string
	}{
		{"gemma4:e4b", "gemma4:e4b", ""},
		{"qwen3:8b?max_output_tokens=8192", "qwen3:8b", "max_output_tokens=8192"},
		{"m?a=1&b=2", "m", "a=1&b=2"},
		{"m?", "m", ""},
	}
	for _, c := range cases {
		t.Run(c.in, func(t *testing.T) {
			model, raw := splitModelQuery(c.in)
			if model != c.wantModel || raw != c.wantRaw {
				t.Fatalf("splitModelQuery(%q) = (%q,%q), want (%q,%q)", c.in, model, raw, c.wantModel, c.wantRaw)
			}
		})
	}
}
