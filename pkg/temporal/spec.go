//ff:type feature=temporal type=model
//ff:what AI가 원문 읽기로 채운 구조화 시간 명세. Kind/Calendar/Start/End/OffsetDays/Anchors를 담고, 정규 Value는 검증기(Resolve)가 산출한다. AI는 식별만, 변환·산술은 기계.

package temporal

// Spec is the structured time description a gate rule fills from an AI's reading. The
// verifier produces the normalized Value; the AI supplies the identification only.
type Spec struct {
	Kind       Kind     `json:"kind"`
	Calendar   Calendar `json:"calendar,omitempty"`
	Start      string   `json:"start,omitempty"`
	End        string   `json:"end,omitempty"`
	OffsetDays int      `json:"offset_days,omitempty"`
	Anchors    []string `json:"anchors,omitempty"`
}
