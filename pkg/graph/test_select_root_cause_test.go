//ff:func feature=graph type=helper control=sequence
//ff:what selectRootCause 단위증명 — 잔존 카운터에서 결정타·곁가지를 가르는 분기를 테이블로 커버한다: ①상류(fmt, reach 큼)가 결정타·나머지 Fail은 siblings, ②상하관계 없는 Fail 둘이면 등록순 첫째가 결정타, ③잔존 Fail 0(Review만)이면 첫 잔존이 결정타·siblings 없음, ④빈 입력이면 빈 결정타.

package graph

import (
	"testing"

	"github.com/park-jun-woo/reins/pkg/gate"
)

func TestSelectRootCause(t *testing.T) {
	g, byID := buildSupersessionGraph()
	byID["fmt"].level = gate.LevelFail
	byID["holder"].level = gate.LevelFail
	byID["free"].level = gate.LevelReview

	ac := func(id string) activeCounter { return activeCounter{node: byID[id]} }

	t.Run("upstream fmt is root, holder is sibling", func(t *testing.T) {
		root, sibs := g.selectRootCause([]activeCounter{ac("fmt"), ac("holder")})
		if root.node.id != "fmt" {
			t.Fatalf("root=%s want fmt", root.node.id)
		}
		if len(sibs) != 1 || sibs[0].node.id != "holder" {
			t.Fatalf("siblings=%v want [holder]", acIDs(sibs))
		}
	})

	t.Run("later Fail with greater reach wins", func(t *testing.T) {
		// holder first (reach 0) but fmt later (reach 2) -> fmt is the root cause.
		root, sibs := g.selectRootCause([]activeCounter{ac("holder"), ac("fmt")})
		if root.node.id != "fmt" {
			t.Fatalf("root=%s want fmt", root.node.id)
		}
		if len(sibs) != 1 || sibs[0].node.id != "holder" {
			t.Fatalf("siblings=%v want [holder]", acIDs(sibs))
		}
	})

	t.Run("no relation -> registration order first", func(t *testing.T) {
		g2 := NewGraph("flat")
		a := mkNode(g2, "a")
		b := mkNode(g2, "b")
		g2.nodes = []*Node{a, b}
		root, sibs := g2.selectRootCause([]activeCounter{{node: a}, {node: b}})
		if root.node.id != "a" {
			t.Fatalf("root=%s want a", root.node.id)
		}
		if len(sibs) != 1 || sibs[0].node.id != "b" {
			t.Fatalf("siblings=%v want [b]", acIDs(sibs))
		}
	})

	t.Run("review-only -> first remaining is focus, no siblings", func(t *testing.T) {
		root, sibs := g.selectRootCause([]activeCounter{ac("free")})
		if root.node.id != "free" {
			t.Fatalf("root=%s want free", root.node.id)
		}
		if len(sibs) != 0 {
			t.Fatalf("siblings=%v want none", acIDs(sibs))
		}
	})

	t.Run("empty -> zero root", func(t *testing.T) {
		root, sibs := g.selectRootCause(nil)
		if root.node != nil || len(sibs) != 0 {
			t.Fatalf("empty: root=%v sibs=%v", root.node, sibs)
		}
	})
}
