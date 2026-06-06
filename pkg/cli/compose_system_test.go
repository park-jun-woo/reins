//ff:func feature=cli type=helper control=iteration dimension=1
//ff:what TestComposeSystem — 세 분기를 테이블로 검증: 빈 전역→fallback, 빈 코칭→전역만, 둘 다 있으면 줄바꿈 결합.

package cli

import (
	"testing"
)

// TestComposeSystem covers the three branches: empty global falls back to the
// generic prompt; empty coaching returns the global unchanged; both present are
// joined with a newline.
func TestComposeSystem(t *testing.T) {
	cases := []struct {
		name   string
		global string
		coach  string
		want   string
	}{
		{"both-empty", "", "", fallbackSystem},
		{"global-only", "G", "", "G"},
		{"fallback-with-coach", "", "C", fallbackSystem + "\nC"},
		{"both", "G", "C", "G\nC"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := composeSystem(c.global, c.coach)
			if got != c.want {
				t.Fatalf("composeSystem(%q,%q) = %q, want %q", c.global, c.coach, got, c.want)
			}
		})
	}
}
