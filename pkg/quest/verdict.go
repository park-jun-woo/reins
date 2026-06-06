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
	// RootCause is the ID of the top rule that produced a FAIL/REVIEW (the first
	// fired Fail rule on the flat path, or the graph backend's selected root-cause
	// counter). It is additive and backward-compatible: empty on PASS or when no
	// backend fills it. The agent loop indexes RuleSystem[RootCause] for rule-specific
	// coaching, giving the flat (ccnews) and graph (comail) paths one code path.
	RootCause string `json:"root_cause,omitempty"`
}
