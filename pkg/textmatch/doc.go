// Package textmatch is reins' deterministic "does this token literally appear in
// the source" verifier — the universal anti-hallucination primitive that quest
// gate rules call for cheese defense (reins Phase003 step ④, Phase005).
//
// It compares surface forms only: both source and token are normalized identically
// (Unicode NFC + whitespace collapse) and then matched as a substring. No
// translation, no synonyms, no fuzzy matching — same (source, token) always yields
// the same answer. The NFC step fixes the Bengali/diacritic false-negatives that a
// whitespace-only normalize misses (composed vs decomposed forms).
package textmatch
