//ff:func feature=graph type=helper control=iteration dimension=1
//ff:what shortNameFor 증명 — toulmin shortName 규칙 재현: base의 마지막 "."까지 잘라내고 "#spec"는 유지한다. spec 유무·dot 유무 4분기를 테이블로 커버한다.

package graph

import "testing"

func TestShortNameFor(t *testing.T) {
	cases := []struct {
		full string
		want string
	}{
		{"pkg/graph.alwaysTrue#warrant#0", "alwaysTrue#warrant#0"}, // dot + spec
		{"alwaysTrue#warrant#0", "alwaysTrue#warrant#0"},           // no dot + spec
		{"a/b/c.fn", "fn"}, // dot, no spec
		{"plain", "plain"}, // no dot, no spec
		{".fn#s", "fn#s"},  // leading dot
		{"a.b.c#x", "c#x"}, // multiple dots -> last
	}
	for _, c := range cases {
		if got := shortNameFor(c.full); got != c.want {
			t.Fatalf("shortNameFor(%q)=%q want %q", c.full, got, c.want)
		}
	}
}
