//ff:type feature=quest type=model
//ff:what 게이트의 제출 1건 판정 타입과 5개 상수. State와 별개 — FAIL은 재시도 가능한 실패 시도다.

package quest

// Outcome is the gate's deterministic judgment of one submission. It is distinct
// from State: FAIL is a retryable failed attempt, not an item end state.
type Outcome string

const (
	OutPass   Outcome = "PASS"
	OutReview Outcome = "REVIEW"
	OutFail   Outcome = "FAIL"
	OutSkip   Outcome = "SKIPPED"
	OutBlock  Outcome = "BLOCKED"
)
