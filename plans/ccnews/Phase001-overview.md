# Phase 001 — ccnews 개요 & 아키텍처

- 프로젝트: ccnews (GitHub: `park-jun-woo/quest-ccnews`)
- 작성일: 2026-06-03
- 상태: 설계 확정 · 8-Phase 분할

## 목적

CC-NEWS 덤프에서 뉴스 기사를 읽어 **육하원칙(event6)** 을 추출하는 **퀘스트 CLI**.
[comail](https://github.com/park-jun-woo/quest-comail)의 무단수집·개인정보 이슈를 피한,
**공개 가능한** 퀘스트 위력 데모.

핵심 명제: *"AI가 요약해도 끝이 아니다. 추출한 사실이 원문에 실재하는지 기계가 대조해
통과한 것만 잠근다. 완료 판정 권한은 AI에게 없다."* (생성은 AI, 판정은 기계)

## 적법성 설계 (왜 공개 가능한가)

| 축 | 근거 |
|---|---|
| 수집 | 본문은 Common Crawl **CC-NEWS 덤프**(합법·공개)에서. 라이브 재크롤 안 함. |
| robots 재확인 | CC 수집 *당시* ≠ *현재* robots일 수 있고 완전보장 아님 → **현재 robots 1회 재확인**, 거부는 BLOCKED(감사). |
| 신뢰 | **구조화 데이터(JSON-LD/OG)를 self-declare한 사이트만** 처리. 없으면 SKIPPED(신뢰 불가). |
| 공개 | **사실만 농축**(사실은 저작권 비보호) + 원문 링크. **본문 텍스트 미저장**(표현 복제 회피). |

## 아키텍처 (2층 퀘스트 + 호스트 캐시)

```
호스트 단위 (캐싱):
  ① robots 게이트       현재 robots.txt 1회 fetch·파싱   ← 유일한 라이브 네트워크
기사 단위:
  host robots에 path 평가 → 거부면 BLOCKED
  구조화 데이터 추출(JSON-LD/OG) → 없으면 SKIPPED(신뢰 불가)
  AI가 event6 추출 → 사실앵커 게이트 → PASS/REVIEW
```

- 추출은 **구조화 데이터만**(셀렉터/AI 학습 없음 — Phase005). AI는 event6 생성에만.
- 앵커는 **원어** 토큰으로 환각 차단, 산출 `value`는 **영어**로 통일(Phase006).

## verdict / 상태 (상세 Phase002·006)

`TODO` → `PASS`(필수 앵커 완료) / `REVIEW`(해석필드 사람확인) / `DONE`(시도 소진) /
`BLOCKED`(robots 거부) / `SKIPPED`(구조화 데이터 없음). 단방향 잠금.

## CLI 표면 (잠정)

| 명령 | 역할 |
|---|---|
| `ccnews run` | 투트랙 인제스천 루프(다운로드→레코드→기사 추가→커서 전진). `--track forward\|backward\|both` |
| `ccnews next` | 다음 TODO 기사 1건 + 원문 본문 + event6 작성 프롬프트 출력 |
| `ccnews submit` | `--url`/`--event6`로 제출된 event6를 앵커 게이트로 검증·잠금/재시도 |
| (export 없음) | 종단 기사는 `--out` JSONL로 **증분 sweep 저장**(Phase007). xlsx/csv 일괄 export는 폐기 |

## 확정된 핵심 결정

1. CC-NEWS **전체**, **투트랙**(forward 최신→대기 / backward 과거 고반복) — Phase003
2. WARC **자동 다운로드**, 최신부터, CLI가 커서로 관리 — Phase003
3. **산출 영어 통일**, `anchors`=원어, 기사 `lang` 명시, 검증은 원어 — Phase006
4. **robots UA** = `parkjunwoo-quest/0.1 (+https://www.parkjunwoo.com)` → 실질적으로 `*` 그룹 적용 — Phase004
5. **구조화 데이터 1단계만 신뢰** — JSON-LD/OG 없으면 SKIPPED — Phase005
6. **출력 = JSONL 증분 append** (xlsx/csv export 폐기). 경로 `--out`, 미지정 시 기본값 — Phase007

## Phase 인덱스

| Phase | 문서 | 범위 |
|---|---|---|
| 001 | `overview` | 본 문서 — 설계 지도 |
| 002 | `session-schema` | session.json 스키마 + 상태기계 + `internal/session` |
| 003 | `ingestion` | 투트랙 WARC 자동 다운로드 + 커서 + scan |
| 004 | `robots-gate` | robots fetch·파서·캐시·path평가 → BLOCKED |
| 005 | `structured-extraction` | JSON-LD/OG 추출 + 신뢰 게이트 → SKIPPED |
| 006 | `event6-anchor-gate` | event6 추출 + 원어앵커/영어산출 + verdict |
| 007 | `output` | 확정 기사 JSONL 증분 저장(`--out`, 기본값) |
| 008 | `charset` | 문자셋 감지 → UTF-8 정규화 (ReadBody 디코드, 비-UTF8 모지바케 해소) |
| 009 | `anchor-hardening` | event6 앵커 게이트 치즈 저항 강화 (빈/흔한 앵커·공허 value 차단) |
| 010 | `when-optional` | when 필수 해제 + 발행일 구조화 데이터 직공급(extract.Apply 배선) |

## 다음 단계
Phase002(세션 스키마) 확정 → 구현은 002→003→…→007 순서(의존 체인). Phase008(charset)은
003(ReadBody) 디코드 보강으로 003 이후 어느 시점에나 합류 — 005/006의 추출·앵커 품질 선결조건.
