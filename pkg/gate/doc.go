//ff:type feature=gate type=model
//ff:what gate 패키지 개요 — 퀘스트 게이트는 위반 탐지 규칙의 카탈로그(yongol/toulmin 패턴). 각 규칙이 문제 발견 시 발동(true)+quest.Fact를 내고, 심각도는 Level(Fail/Review). Evaluate가 카탈로그를 1건에 평가해 레벨로 집계(가중치 아님). defeat/h-Categoriser는 toulmin 백엔드 도입 시 예약된 미래 플러그인.

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
