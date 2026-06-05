//ff:func feature=graph type=helper control=iteration dimension=1
//ff:what alwaysTrue 증명 — 항상 활성 워런트 fn이 어떤 ctx/specs에서도 (true, nil)을 내는지 검증한다.

package graph

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestAlwaysTrue(t *testing.T) {
	cases := []toulmin.Context{
		nil,
		toulmin.NewContext(),
	}
	for i, ctx := range cases {
		ok, ev := alwaysTrue(ctx, nil)
		if !ok || ev != nil {
			t.Fatalf("case %d: got (%v,%v) want (true,nil)", i, ok, ev)
		}
	}
}
