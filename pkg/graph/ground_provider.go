//ff:type feature=graph type=model
//ff:what GroundProvider — ground 이름(예 "source-body")을 실제 값으로 푸는 소비자 공급 함수. ctx(제출 사실에서 URL/도메인 추출)와 요청당 스냅샷(snap.HTTPBody/MXResolves로 lazy·캐시 resolve)을 받아 (값, 에러)를 낸다. reins는 어떤 이름이 어떤 URL인지 모르므로(generic) 소비자가 매핑을 공급한다. 에러는 staged 평가가 tier0 통과 후 terminal FAIL Fact로 환원(열린결정 #8).

package graph

import (
	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/ground"
)

// GroundProvider is the consumer-supplied function that resolves a ground name
// (e.g. "source-body") to its value. It receives the gate.Context (to extract the
// URL/domain from the submission facts) and the per-request ground.Snapshot (whose
// HTTPBody/MXResolves resolve lazily and cache), and returns the resolved value or
// an error. reins is generic — it does not know which name maps to which URL — so
// the consumer supplies the mapping. The staged evaluator reduces a non-nil error to
// a terminal FAIL Fact after tier-0 has passed (open decision #8).
type GroundProvider func(name string, ctx gate.Context, snap *ground.Snapshot) (string, error)
