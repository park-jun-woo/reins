//ff:func feature=graph type=helper control=sequence
//ff:what collectSupersessionEvents 단위증명 — 활성 fmt가 활성 holder·free를 흡수하면 2 이벤트(absorbed=holder/free, by=fmt)를 등록순으로 내는지, 하류가 비활성이면 이벤트가 없는지, 같은 하류를 흡수하는 상류가 둘이면 등록순 첫 상류가 by인지(fmt before free) 검증한다.

package graph

import "testing"

func TestCollectSupersessionEvents(t *testing.T) {
	g, byID := buildSupersessionGraph()
	ac := func(id string) activeCounter { return activeCounter{node: byID[id]} }

	t.Run("fmt absorbs active holder+free", func(t *testing.T) {
		evs := g.collectSupersessionEvents([]activeCounter{ac("fmt"), ac("holder"), ac("free")})
		if len(evs) != 2 {
			t.Fatalf("events=%d want 2", len(evs))
		}
		// registration order: holder before free; both superseded by the first upstream (fmt).
		if evs[0].absorbed.node.id != "holder" || evs[0].by.node.id != "fmt" {
			t.Fatalf("ev0 = %s by %s want holder by fmt", evs[0].absorbed.node.id, evs[0].by.node.id)
		}
		if evs[1].absorbed.node.id != "free" || evs[1].by.node.id != "fmt" {
			t.Fatalf("ev1 = %s by %s want free by fmt", evs[1].absorbed.node.id, evs[1].by.node.id)
		}
	})

	t.Run("downstream inactive -> no event", func(t *testing.T) {
		evs := g.collectSupersessionEvents([]activeCounter{ac("fmt")})
		if len(evs) != 0 {
			t.Fatalf("events=%d want 0 (downstream not active)", len(evs))
		}
	})

	t.Run("free-only absorbs holder", func(t *testing.T) {
		evs := g.collectSupersessionEvents([]activeCounter{ac("holder"), ac("free")})
		if len(evs) != 1 || evs[0].absorbed.node.id != "holder" || evs[0].by.node.id != "free" {
			t.Fatalf("events=%v want [holder by free]", evs)
		}
	})
}
