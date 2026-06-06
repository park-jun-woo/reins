//ff:func feature=llm type=helper control=iteration dimension=1
//ff:what TestFromFlagErrors — 실패 케이스(콜론 없음·빈 model·지원 외 backend·빈 flag)를 테이블로 순회하며 에러+nil backend를 검증.

package llm

import (
	"testing"
)

// TestFromFlagErrors covers the failure modes: no colon, empty model, unknown backend.
func TestFromFlagErrors(t *testing.T) {
	cases := []string{
		"ollama",      // no colon
		"ollama:",     // empty model
		"unknown:foo", // unsupported backend
		"",            // empty flag (no colon)
	}
	for _, flag := range cases {
		t.Run(flag, func(t *testing.T) {
			b, err := FromFlag(flag)
			if err == nil {
				t.Fatalf("FromFlag(%q) = %v, want error", flag, b)
			}
			if b != nil {
				t.Fatalf("FromFlag(%q) backend = %v, want nil on error", flag, b)
			}
		})
	}
}
