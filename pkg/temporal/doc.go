//ff:type feature=temporal type=model
//ff:what temporal 패키지 개요 — AI가 식별한 구조화 시간 명세(역법/성분/오프셋)를 정규 그레고리력 ISO(단일·기간)로 결정론 변환·검증하는 순수 헬퍼. v1=그레고리+상대+기간, 비그레고리·미파싱은 Determined=false(REVIEW). 역법 변환·비라틴 숫자표는 v2.

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
