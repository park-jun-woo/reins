// Package quest is reins' irreversible-progress core: the one-way ratchet state
// machine, session persistence, progress tally, verdict application, and terminal
// export. It is pure (no Cobra, no toulmin, no domain types) — "the agent is
// disposable, progress accumulates": once an item reaches a terminal state it never
// reopens, so remaining(t+1) <= remaining(t).
//
// quest does not judge. It applies a gate Verdict (see package gate) to the ratchet
// and carries the Facts back. PASS/REVIEW/SKIPPED/BLOCKED lock an item; FAIL is a
// failed attempt (tries++, → DONE at MaxTries).
package quest
