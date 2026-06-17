//ff:func feature=llm type=helper control=iteration dimension=1 level=error
//ff:what TestCheckOptsAllowed — present⊆allowed면 통과, 벗어난 키가 있으면 에러, 빈 present는 통과, 빈 allowed(subprocess)+키 있으면 에러를 테이블로 검증.

package llm

import "testing"

// TestCheckOptsAllowed: present must be a subset of allowed; an empty allowed set
// rejects any present key; an empty present always passes.
func TestCheckOptsAllowed(t *testing.T) {
	cases := []struct {
		name    string
		present map[string]bool
		allowed []string
		wantErr bool
	}{
		{"subset-ok", map[string]bool{"temperature": true}, []string{"max_output_tokens", "temperature"}, false},
		{"disallowed-key", map[string]bool{"num_ctx": true}, []string{"max_output_tokens", "temperature"}, true},
		{"empty-present-ok", map[string]bool{}, []string{"max_output_tokens"}, false},
		{"subprocess-rejects", map[string]bool{"temperature": true}, nil, true},
		{"empty-both-ok", map[string]bool{}, nil, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := checkOptsAllowed(c.present, c.allowed...)
			if (err != nil) != c.wantErr {
				t.Fatalf("checkOptsAllowed err = %v, wantErr %v", err, c.wantErr)
			}
		})
	}
}
