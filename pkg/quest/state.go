//ff:type feature=quest type=model
//ff:what 퀘스트 아이템의 단방향 래칫 상태 타입과 6개 상수. terminal 상태는 잠금이라 다시 전이하지 않는다.

package quest

// State is the one-way ratchet state of a quest item. A terminal state is a lock:
// it never transitions again.
type State string

const (
	TODO    State = "TODO"
	PASS    State = "PASS"
	REVIEW  State = "REVIEW"
	DONE    State = "DONE"
	SKIPPED State = "SKIPPED"
	BLOCKED State = "BLOCKED"
)

// Terminal reports whether s is a locked end state (no further transitions).
func (s State) Terminal() bool {
	switch s {
	case PASS, REVIEW, DONE, SKIPPED, BLOCKED:
		return true
	}
	return false
}
