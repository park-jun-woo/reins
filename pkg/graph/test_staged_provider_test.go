//ff:func feature=graph type=helper control=sequence
//ff:what 테스트 헬퍼 — "source-body" ground를 고정 URL의 snapshot.HTTPBody로 resolve하는 provider. 두 tier-1 ground가 같은 URL로 매핑되므로 snapshot이 캐시해야 한다.

package graph

import (
	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/ground"
)

// stagedProvider resolves "source-body" by reading the snapshot's HTTPBody for a
// fixed URL. Both tier-1 grounds map to the same URL → snapshot must cache.
func stagedProvider(name string, _ gate.Context, snap *ground.Snapshot) (string, error) {
	return snap.HTTPBody("http://src")
}
