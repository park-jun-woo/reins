// Package gate is reins' deterministic verifier framework: a quest gate is a
// catalog of violation-detecting rules (the yongol/toulmin pattern). Each rule fires
// (true) when it finds a problem and emits a quest.Fact; severity is the rule's
// Level (Fail/Review). Evaluate runs the catalog over one submission and aggregates
// by level — any Fail-rule → FAIL, else any Review-rule → REVIEW, else PASS.
//
// Severity is expressed as a Level, never as a weight, so a single decisive
// violation is FAIL (not a weighted score that nets to REVIEW). A defeasible
// weighting backend (toulmin h-Categoriser) is reserved for genuine competing
// evidence such as L2 AI consensus; the core path here is level aggregation.
//
// A quest plugs in a Definition (Seed/Render/Prepare/Rules); reins drives the rest.
// Rule metadata (ID/Level/Desc) is the auto-generated rulebook — the grep-able
// catalog of "every cheese this gate blocks" (see the cli `rules` command).
package gate
