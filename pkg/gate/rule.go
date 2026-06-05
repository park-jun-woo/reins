package gate

import "github.com/park-jun-woo/reins/pkg/quest"

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

// RuleMeta is the self-documenting catalog entry for a rule (the auto rulebook).
type RuleMeta struct {
	ID    string `json:"id"`
	Level Level  `json:"level"`
	Desc  string `json:"desc"`
}

// Context carries the per-submission facts a rule inspects. Source is the cached
// ground truth (re-confirmed by cheese-defense rules); Submission is the decoded
// domain artifact.
type Context struct {
	Item       *quest.Item
	Submission any
	Source     string
}

// Rule is a violation detector. Check returns fired=true with a Fact when it finds a
// problem; no fire means no problem. The aggregate of fired rules' levels decides
// the verdict (see Evaluate).
type Rule struct {
	Meta  RuleMeta
	Check func(ctx Context) (fired bool, fact quest.Fact)
}

// Definition is the per-quest domain contract. Implement these four and reins
// supplies the ratchet, command skeleton, aggregation, and export.
type Definition interface {
	// Seed creates the initial TODO items from CLI args (files, dirs, a stream).
	Seed(args []string) ([]*quest.Item, error)
	// Render returns the authoring prompt + verification context shown by `next`.
	Render(it *quest.Item) (string, error)
	// Prepare decodes a raw submission into an evaluation Context. A non-nil short
	// verdict short-circuits the gate (e.g. SKIPPED when the source is untrusted).
	Prepare(it *quest.Item, raw []byte) (ctx Context, short *quest.Verdict, err error)
	// Rules is the gate's violation-rule catalog.
	Rules() []Rule
}
