//ff:func feature=llm type=helper control=iteration dimension=1
//ff:what TestClampCtxLo — 하한 미만 step은 2048로 올리고, 하한 이상은 그대로 두는지 검증.

package llm

import "testing"

// TestClampCtxLo: a step below ctxLo is raised to 2048; a step at or above it is
// returned unchanged.
func TestClampCtxLo(t *testing.T) {
	cases := []struct {
		name string
		in   int
		want int
	}{
		{"below lower bound", 1024, 2048},
		{"zero", 0, 2048},
		{"at lower bound", 2048, 2048},
		{"above lower bound", 8192, 8192},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := clampCtxLo(c.in); got != c.want {
				t.Fatalf("clampCtxLo(%d) = %d, want %d", c.in, got, c.want)
			}
		})
	}
}
