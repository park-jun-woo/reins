# Phase 007 — 출력 (JSONL 증분 저장)

- 프로젝트: ccnews · 작성일: 2026-06-03 · 상태: 설계 확정
- 의존: Phase002(세션), Phase006(verdict)

## 목적

확정된 기사를 **JSONL(한 줄당 한 건)** 로 **증분 append** 한다. xlsx/csv 일괄 export는
폐기(대규모·무기한 투트랙이라 "끝나고 한 번에 저장"하는 시점이 없고, 거대 파일 통째 재기록은
비현실적). 경로는 `--out`으로 지정, 미지정 시 기본값.

## 왜 JSONL append인가
- 인제스천은 무기한 → 종료-시점-export가 없음. 기사가 **종단상태(PASS/REVIEW/DONE/BLOCKED/SKIPPED)에
  도달할 때마다 한 줄 append**.
- JSON 배열은 저장 시 전체 재기록(수GB 매번) → 불가. **JSONL은 줄 추가만(O(1))**.
- 부수효과: `session.json`은 작업상태(커서·robots 캐시·진행 중 TODO)만 유지 → 비대화 방지.
  (확정분은 JSONL로 흘려보냄. 세션 경계 유지)

## 출력 레코드 (한 줄 = 기사 1건)

`status`로 collected(PASS/REVIEW)와 audit(DONE/BLOCKED/SKIPPED)를 한 파일에서 구분 — 경로 하나.

```jsonc
// PASS/REVIEW — 사용자 명시 수집항목 전부 포함
{ "url":"...", "host":"...", "media_name":"...", "site_url":"...",
  "lang":"ko", "published_at":"2026-06-02", "collected_at":"...",
  "status":"PASS",                       // PASS | REVIEW
  "crawl_allowed":true,                  // 크롤링허용여부(사용자 명시)
  "license":null,                        // CC저작권종류(있으면)
  "anchor_summary":"4/4",                // 검증 증거(데모 임팩트)
  "event6": { "who":{"value":"...(en)","anchors":["...(원어)"]},
              "when":{...}, "what":{...}, "where":null, "how":null, "why":null } }

// DONE/BLOCKED/SKIPPED — 제외 감사(legality·소진 증빙)
{ "url":"...", "host":"...", "status":"SKIPPED",
  "reason":"구조화 데이터 없음(JSON-LD/OG) — 신뢰 불가", "crawl_allowed":true }
{ "url":"...", "host":"...", "status":"BLOCKED",
  "reason":"robots Disallow: /premium/", "crawl_allowed":false }
{ "url":"...", "host":"...", "status":"DONE",
  "reason":"...(verdict_reason 폴백 — MaxTries 소진)", "crawl_allowed":true }
```

- DONE(시도 소진)도 audit 레코드로 emit한다(사용자 결정). audit `reason`은 `skip_reason`,
  DONE처럼 `skip_reason`이 없으면 `verdict_reason`으로 폴백(코드 `renderAudit` 반영).

- event6 `value`는 영어, `anchors`는 원어(Phase006). 본문 텍스트는 출력 안 함(사실만).
- 매체명+원문URL = CC-BY 출처표시 충족. collected(PASS/REVIEW)만 추리려면 `status` 필터,
  audit(DONE/BLOCKED/SKIPPED)는 제외 사유 서사용.

## CLI / 동작
- 별도 export 명령 없음. `ccnews run`/`submit`이 종단 도달 시 `--out` 파일에 한 줄 append.
- `--out <path>` 미지정 시 기본값(예: `ccnews-results.jsonl`).
- 같은 url 중복 append 방지: 종단상태 1회만 emit(세션이 emit 여부 추적).

## 결정론적 성격
- 출력은 게이트 확정 결과의 직렬화일 뿐 — 판정/생성 없음(순수 기록).

## 열린 결정
- 기본 출력 경로 파일명 확정(`ccnews-results.jsonl` 등).
- collected/audit를 한 파일(status 필드) vs 두 파일(`--out`/`--audit-out`)로 — 1차: 한 파일.
- `anchor_summary` 노출 형식(예: `4/4`, 또는 필드별 anchored 맵).

## 다음 단계
구현 완료 후: 공개 데모(JSONL 통계 + collected 샘플 + audit 줄로 적법성 서사).
