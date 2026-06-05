//ff:type feature=quest type=model
//ff:what 게이트 결과 — Outcome와 그 근거 Facts(FAIL/REVIEW용)를 담는다.

package quest

// Verdict is a gate result: the outcome plus the Facts behind it (for FAIL/REVIEW).
type Verdict struct {
	Outcome Outcome
	Facts   []Fact
}
