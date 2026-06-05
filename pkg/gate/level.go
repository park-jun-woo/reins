//ff:type feature=gate type=model
//ff:what 규칙의 심각도 레벨. LevelFail(발동 시 FAIL — 완료를 막는 위반), LevelReview(발동 시 REVIEW — 모호, 사람 필요). 가중치가 아니라 레벨이다. String()은 "FAIL"|"REVIEW".

package gate

// Level is a rule's severity: a fired Fail rule makes the gate FAIL, a fired Review
// rule makes it REVIEW (when nothing failed). Severity is a level, not a weight.
type Level int

const (
	LevelFail   Level = iota // fired ⇒ FAIL (a violation that blocks completion)
	LevelReview              // fired ⇒ REVIEW (ambiguous; needs a human)
)

func (l Level) String() string {
	if l == LevelReview {
		return "REVIEW"
	}
	return "FAIL"
}
