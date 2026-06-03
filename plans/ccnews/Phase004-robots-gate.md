# Phase 004 — robots 게이트

- 프로젝트: ccnews · 작성일: 2026-06-03 · 상태: 설계 확정
- 의존: Phase002(세션), Phase003(인제스천)

## 목적

기사 본문은 CC-NEWS 덤프에서 오지만, **현재 robots.txt를 1회 재확인**해 자동접근을 거부한
호스트/경로는 제외한다(보수적·good-faith). 라이브 네트워크는 **호스트당 robots.txt 한 개**뿐.

## 동작

```
호스트 첫 등장 → robots.txt fetch (UA 헤더 = parkjunwoo-quest/...) → 파싱 → hosts[host].robots 캐싱
기사마다 → 캐시된 룰셋에 기사 path 평가(네트워크 0)
   허용 → 다음 게이트(Phase005)
   거부 → state=BLOCKED, skip_reason="robots Disallow: <rule>"
```

## robots 규칙
- **UA 매칭**: `parkjunwoo-quest`. 우릴 콕 집은 사이트는 없으니 실질적으로 **`*` 그룹** 적용.
  (사용자 확정 = UA 문자열 + `*` 그룹 평가. `CCBot` 등 다른 봇 그룹 추가 존중은 **보류/옵션** —
  사용자 미지시이며, CC-NEWS는 덤프이므로 우리는 CCBot으로 행세하지 않는다.)
- **파서**: RFC 9309 수준 — Allow/Disallow **최장 일치** 우선, `*`/`$` 와일드카드, 대소문자, 경로 정규화.
- **robots 없음/4xx** = 허용(status=`missing`). **5xx/타임아웃** = `unreachable` → 보수적으로 그 호스트 보류(재시도 or BLOCKED, 열린결정).
- `crawl_delay`는 기록(라이브 fetch가 robots.txt 하나뿐이라 실효는 작음).

## 결정론적 게이트
- 입력: host, path, 캐시된 룰셋, UA. 출력: allow/deny + 매칭 룰(감사). 의견 없음.
- BLOCKED 기사는 JSONL audit 레코드(`status:"BLOCKED"`, Phase007)로 "robots 존중" 증빙.

## 열린 결정
- robots `unreachable`(5xx/타임아웃) 시 정책: 보류 후 재시도 N회 → 그래도면 BLOCKED vs 허용.
- robots 캐시 TTL(런 단위 캐시면 충분 / 장기 run이면 재확인 주기).

## 다음 단계
Phase005(구조화 추출) — 허용 호스트 기사의 본문 추출 + 신뢰 게이트.
