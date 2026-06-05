//ff:type feature=quest type=model
//ff:what 게이트 규칙이 발동 시 싣는 정량·위치지정 evidence. FAIL 시 에이전트에 그대로 피드백된다.

package quest

// Fact is fact-based feedback a gate rule emits when it fires: located and
// quantified, with "no room to flatter" (how-make-quest). On FAIL it is handed back
// to the agent so a sycophantic model converges on a correction instead of arguing.
type Fact struct {
	Rule     string `json:"rule,omitempty"`
	Where    string `json:"where,omitempty"`
	Expected string `json:"expected,omitempty"`
	Actual   string `json:"actual,omitempty"`
}
