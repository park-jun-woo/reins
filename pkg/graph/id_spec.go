//ff:type feature=graph type=model
//ff:what idSpec — 노드별 고유 toulmin ruleID를 강제하는 식별 Spec. 어댑터 클로저는 같은 소스 위치에서 만들어져 funcID가 충돌하므로(toulmin ruleID = funcID), 노드마다 고유 ID를 Spec으로 실어 ruleID를 분리한다. SpecName/String은 ID 반환, Validate는 빈 ID 거부(별도 파일).

package graph

// idSpec forces a unique toulmin ruleID per node. Adapter closures created at the
// same source location share a funcID, so without a distinguishing spec toulmin's
// ruleID (funcID#spec) would collide and panic on duplicate registration. Carrying
// the node's unique ID as a spec disambiguates each rule.
type idSpec struct {
	ID string
}
