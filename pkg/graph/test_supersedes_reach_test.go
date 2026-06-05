//ff:func feature=graph type=helper control=iteration dimension=1
//ff:what supersedesReach 단위증명 — fmt(holder+free 직접, free 경유 transitive)·free(holder)·holder(0)의 transitive 흡수 수를 테이블로 검증하고, 미등록 id가 0인지(누구도 안 흡수) 확인한다.

package graph

import "testing"

func TestSupersedesReach(t *testing.T) {
	g, _ := buildSupersessionGraph()
	cases := []struct {
		id   string
		want int
	}{
		{"fmt", 2},    // holder + free (free->holder transitive, deduped)
		{"free", 1},   // holder
		{"holder", 0}, // absorbs nothing
		{"absent", 0}, // unknown id absorbs nothing
	}
	for _, c := range cases {
		if got := g.supersedesReach(c.id); got != c.want {
			t.Errorf("supersedesReach(%q) = %d, want %d", c.id, got, c.want)
		}
	}
}
