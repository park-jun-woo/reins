//ff:type feature=ground type=model
//ff:what ground 패키지 개요 — 검증 "ground"(외부 사실) 원시연산. textmatch/temporal이 글자·시간 원시연산이듯 네트워크 ground도 reins 자산: HTTPBody(URL→본문)·MXResolves(도메인→수신가능). 핵심 불변식: 평가 1회당 1회 스냅샷·캐시(요청당 결정성) + resolver 주입(테스트는 실제 네트워크 금지). resolve 에러는 값이 아니라 호출측이 결정론 FAIL Fact로 환원할 신호.

// Package ground provides reins' verification "ground" primitives — the external
// facts a gate rule needs but cannot compute purely. Just as textmatch/temporal are
// reins' character/time primitives, the network ground is a reins asset:
// HTTPBody(url) → page body and MXResolves(domain) → deliverable bool. Consumers do
// not hand-code fetch/MX; they read ground.
//
// The load-bearing invariant is per-request snapshot+cache: within one evaluation a
// ground value is resolved at most once (on first read) and cached, so every leaf
// the rule reads is fixed for that evaluation (determinism; retry happens on the
// next submit). The resolver is injectable — the default uses the real net stack
// (http.Get / net.LookupMX), and tests inject a fake so no real network is touched.
// A resolve error is surfaced, not swallowed; the caller reduces it to a
// deterministic FAIL Fact (Phase007 §④/§⑤, open decision #8).
package ground
