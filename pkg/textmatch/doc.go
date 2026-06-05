//ff:type feature=textmatch type=model
//ff:what textmatch 패키지 개요 — 토큰이 원천에 글자 그대로 substring으로 실재하는지 결정론으로 판정하는 순수 검증기. 원천·토큰 양쪽에 NFC 정규화+공백 접기를 동일 적용한 표면형 비교만 한다(번역·동의어·퍼지 없음). 게이트 규칙 패밀리 ④(치즈 방어)가 호출하는 보편 헬퍼.

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
