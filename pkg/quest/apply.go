//ff:func feature=quest type=helper control=selection
//ff:what verdict를 래칫에 적용한다. PASS/REVIEW/SKIPPED/BLOCKED는 잠금, FAIL은 Tries++(MaxTries면 DONE 잠금). now 주입(순수).

package quest

// Apply transitions an item per a gate Verdict. It is a pure state mutation (now is
// injected). PASS/REVIEW/SKIPPED/BLOCKED lock the item and stamp CollectedAt; FAIL
// is a failed attempt (Tries++, locking to DONE once Tries reaches MaxTries). Every
// call appends an Attempt to the log.
func Apply(it *Item, v Verdict, now string) {
	it.Log = append(it.Log, Attempt{Try: it.Tries + 1, Outcome: string(v.Outcome), Reason: v.Reason()})
	switch v.Outcome {
	case OutPass:
		lock(it, PASS, v, now)
	case OutReview:
		lock(it, REVIEW, v, now)
	case OutSkip:
		lock(it, SKIPPED, v, now)
	case OutBlock:
		lock(it, BLOCKED, v, now)
	default: // OutFail
		it.Tries++
		if it.Tries >= MaxTries {
			lock(it, DONE, v, now)
		}
	}
}
