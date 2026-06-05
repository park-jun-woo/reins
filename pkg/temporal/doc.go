// Package temporal is reins' deterministic time-spec normalizer (a gate rule helper,
// reins Phase006). It takes a structured time Spec — calendar/components/offset that
// an AI identified by *reading* the source — and converts it to a canonical
// Gregorian ISO value (single date or interval). Multilingual reading stays with the
// AI; this package only does the (language-independent, deterministic) conversion and
// arithmetic.
//
// v1 scope is Gregorian + interval + relative (ref-based) — it covers the bulk of
// news with only the standard library. Any non-Gregorian calendar, an unparseable
// component, or a relative spec without a ref returns Determined=false (mapped to a
// REVIEW-level rule), honest about what it cannot yet decide. v2 (deferred) would add
// calendar-conversion and extra-numeral-system tables here; reins does not pull in
// those libraries until the need is observed.
package temporal
