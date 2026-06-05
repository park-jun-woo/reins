//ff:func feature=graph type=helper control=iteration dimension=1
//ff:what counterName 단위증명 — Fact.Rule이 있으면 그 값을, 없으면 노드 shortName으로 폴백하는지 두 케이스로 검증한다.

package graph

import (
	"testing"

	"github.com/park-jun-woo/reins/pkg/quest"
)

func TestCounterName(t *testing.T) {
	g, byID := buildSupersessionGraph()
	cases := []struct {
		name string
		ac   activeCounter
		want string
	}{
		{"fact rule wins", activeCounter{node: byID["fmt"], fact: quest.Fact{Rule: "email-format"}}, "email-format"},
		{"shortName fallback", activeCounter{node: byID["holder"]}, "holder"},
	}
	_ = g
	for _, c := range cases {
		if got := counterName(c.ac); got != c.want {
			t.Errorf("%s: counterName = %q, want %q", c.name, got, c.want)
		}
	}
}
