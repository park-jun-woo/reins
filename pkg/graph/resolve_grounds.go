//ff:func feature=graph type=helper control=iteration dimension=1
//ff:what resolveGrounds — tier1이 선언한 ground 이름들을 provider로 1회씩 lazy resolve해 ctx.Grounds에 채운다(스냅샷 캐시로 같은 ground 재읽기는 resolver 0회 추가). 이름이 없으면 resolve·provider 호출 0회(네트워크 0). resolve 에러는 즉시 terminal FAIL Verdict(결정론 Fact)로 반환하고 그 자리에서 멈춘다(열린결정 #8). 성공 시 grounds 맵과 nil Verdict 포인터를 낸다.

package graph

import (
	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/ground"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// resolveGrounds resolves each ground name declared by tier-1 counters once via the
// provider and fills the returned grounds map (the snapshot caches, so a repeat read
// of the same ground adds zero resolver calls). With no names it makes zero provider
// calls (zero network). A resolve error is returned immediately as a terminal FAIL
// Verdict carrying a deterministic Fact, stopping at that name (open decision #8). On
// success it returns the grounds map and a nil Verdict pointer.
func (g *Graph) resolveGrounds(ctx gate.Context, snap *ground.Snapshot, provider GroundProvider) (map[string]string, *quest.Verdict) {
	names := g.tier1GroundNames()
	if len(names) == 0 {
		return nil, nil
	}
	grounds := make(map[string]string, len(names))
	for _, name := range names {
		val, err := provider(name, ctx, snap)
		if err != nil {
			v := quest.Verdict{
				Outcome: quest.OutFail,
				Facts: []quest.Fact{{
					Rule:   "ground",
					Where:  name,
					Actual: err.Error(),
				}},
			}
			return nil, &v
		}
		grounds[name] = val
	}
	return grounds, nil
}
