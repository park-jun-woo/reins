package quest

import "strings"

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

// Verdict is a gate result: the outcome plus the Facts behind it (for FAIL/REVIEW).
type Verdict struct {
	Outcome Outcome
	Facts   []Fact
}

// Reason renders the Facts into a single human-readable line for the audit log.
func (v Verdict) Reason() string {
	if len(v.Facts) == 0 {
		return string(v.Outcome)
	}
	parts := make([]string, 0, len(v.Facts))
	for _, f := range v.Facts {
		seg := f.Where
		if f.Expected != "" || f.Actual != "" {
			seg += " (expected " + f.Expected + ", got " + f.Actual + ")"
		}
		if f.Rule != "" {
			seg = f.Rule + ": " + seg
		}
		parts = append(parts, seg)
	}
	return strings.Join(parts, "; ")
}
