//ff:func feature=graph type=helper control=iteration dimension=1
//ff:what SpecName 증명 — idSpec.SpecName이 ID 문자열을 그대로 반환하는지 검증한다.

package graph

import "testing"

func TestIDSpecSpecName(t *testing.T) {
	cases := []string{"", "warrant#0", "rule.dup#n3"}
	for _, id := range cases {
		if got := (idSpec{ID: id}).SpecName(); got != id {
			t.Fatalf("SpecName()=%q want %q", got, id)
		}
	}
}
