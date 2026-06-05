//ff:func feature=graph type=helper control=sequence
//ff:what cloneNode 워런트 경로 증명 — 워런트 노드를 dst 그래프에 복제하면 isWarrant·id·shortName이 보존되고 graph가 dst로 재바인딩되며 toulmin rule이 채워지고 dst.nodes에 추가되는지 검증한다(네트워크 0, 결정적).

package graph

import (
	"testing"
)

func TestCloneNodeWarrant(t *testing.T) {
	src := NewGraph("src")
	w := src.Warrant(alwaysTrue) // registers a warrant node

	dst := NewGraph("dst")
	got := dst.cloneNode(w)

	if !got.isWarrant {
		t.Fatalf("cloned warrant: isWarrant=false want true")
	}
	if got.id != w.id || got.shortName != w.shortName {
		t.Fatalf("cloned warrant: id=%q short=%q want id=%q short=%q",
			got.id, got.shortName, w.id, w.shortName)
	}
	if got.graph != dst {
		t.Fatalf("cloned warrant: graph not rebound to dst")
	}
	if got.tr == nil {
		t.Fatalf("cloned warrant: toulmin rule is nil")
	}
	if len(dst.nodes) != 1 || dst.nodes[0] != got {
		t.Fatalf("cloned warrant: dst.nodes=%v want [got]", dst.nodes)
	}
}
