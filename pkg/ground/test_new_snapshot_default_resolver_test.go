//ff:func feature=ground type=helper control=sequence
//ff:what NewSnapshot(nil) 기본 resolver 채택 증명 — resolver를 넘기지 않으면 non-nil 기본(실네트워크) resolver가 들어가는지 검증한다. 여기서 네트워크는 건드리지 않는다.

package ground

import "testing"

// TestNewSnapshotDefaultResolver proves NewSnapshot(nil) adopts a non-nil default
// resolver (the real-network one) without touching the network here.
func TestNewSnapshotDefaultResolver(t *testing.T) {
	s := NewSnapshot(nil)
	if s.resolver == nil {
		t.Fatal("want a default resolver, got nil")
	}
}
