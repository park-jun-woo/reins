//ff:type feature=quest type=model
//ff:what 게이트 결과 — Outcome와 그 근거 Facts(FAIL/REVIEW용)를 담는다.

package quest

// Verdict is a gate result: the outcome plus the Facts behind it (for FAIL/REVIEW).
// Feedback is an optional pre-rendered "walkthrough" (the graph backend's root-cause
// strategy text). It is additive: empty when a flat rule catalog produced the
// verdict, so Apply/Reason/gate.Evaluate are unaffected. When set (graph FAIL/REVIEW),
// the CLI prints it in place of the flat Fact loop.
type Verdict struct {
	Outcome  Outcome
	Facts    []Fact
	Feedback string `json:"feedback,omitempty"`
}
