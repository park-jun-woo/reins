//ff:func feature=quest type=helper control=sequence
//ff:what 아이템을 terminal 상태로 잠그고 마지막 verdict·reason·수집 시각을 기록한다(Apply 내부용).

package quest

// lock sets the item to a terminal state and records the last verdict, reason, and
// collection time. now is injected so the mutation stays pure.
func lock(it *Item, s State, v Verdict, now string) {
	it.State = s
	it.Verdict = string(v.Outcome)
	it.Reason = v.Reason()
	it.CollectedAt = now
}
