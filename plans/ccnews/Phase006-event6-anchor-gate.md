# Phase 006 — event6 사실 앵커 게이트 & verdict

- 프로젝트: ccnews · 작성일: 2026-06-03 · 상태: 설계 확정
- 의존: Phase002(세션), Phase005(추출)

## 목적

추출 본문에서 **육하원칙(event6)** 을 뽑되, AI가 지어낸 사실을 막는다. 각 사실 토큰이
**원문에 글자 그대로 존재**하는지 기계가 대조해 통과한 것만 잠근다. (생성=AI, 판정=기계)

## 필드 (존재하는 것만)

| 필드 | 필수? | 비고 |
|---|---|---|
| 누가 who / 언제 when / 무엇 what | **필수** | 값+앵커 필수. 비면 FAIL("부재" 변명 불가) |
| 어디서 where / 어떻게 how / 왜 why | 선택 | 있으면 앵커, 없으면 빈칸 허용 |

부재 cheese 차단: 뉴스에 누가·언제·무엇이 없을 리 없으므로 필수 3은 "없음" 불가.

## 산출 언어 = 영어, 앵커 = 원어
- `value` = **영어** 산출 (예: `"Lee Jae-yong, Samsung Chairman"`)
- `anchors` = **원어** 사실 토큰 (예: `["이재용","삼성전자"]`)
- 검증 순서: ① **원어 앵커로 환각 차단** → ② 통과 후 **영어 번역**(게이트 이후).
- 날짜·숫자는 ISO/숫자 정규화(언어 무관). 기사 `lang`에 원문 언어코드.

## 앵커 게이트 (결정론적)
- 각 present 필드의 `anchors[*]`가 **원문 텍스트에 substring 존재**해야 `anchored=true`.
  앵커는 **반드시 원문에 나타난 그대로의 표면형**이어야 한다(예: 원문이 "6월 2일"이면 앵커도 "6월 2일").
- `value`(영어/ISO 정규화)는 **표시용 변환일 뿐 게이트 대상이 아니다**. 게이트는 오직
  `anchors` 표면형 ↔ 원문 substring만 본다. 따라서 `value="2026-06-02"`와 앵커 `"6월 2일"`은
  서로 비교하지 않는다(불가능한 매핑 검증을 시도하지 않음). value의 ISO/영어 정확성은 검증 대상이 아니라 산출물.
- 문장 통째 인용 아님 — 사실 토큰(인명·날짜·숫자·지명) 단위 → 환각 차단 + 표현 저작권 회피.
- 해석필드(how/why)도 동일 규칙: paraphrase한 value 자체는 게이트 안 됨. 그 안의 사실 토큰을
  `anchors`로 분리해 원문 substring 검증. 사실 토큰을 못 뽑으면 `anchors=[]` → 아래 REVIEW.

## 제출 판정 → 상태 전이

게이트는 **1회 제출**을 평가해 셋 중 하나를 낸다. 이 중 PASS/REVIEW만 종단 article state이고,
**FAIL은 article state가 아니라 "이번 시도 실패" 결과**다 — 재시도하며, tries 소진 시 종단 상태 `DONE`이 된다.
(article state 집합은 Phase002와 동일: TODO/PASS/REVIEW/DONE/BLOCKED/SKIPPED — FAIL은 없음.)

| 제출 결과 | 조건 (전부 기계 판정 — 주관 없음) | 전이 |
|---|---|---|
| **PASS** | 필수 who·when·what 전부 value 있음 + anchors 비어있지 않고 전부 원문 substring. present한 선택필드도 동일. | state=PASS(잠금) |
| **REVIEW** | 필수 3은 PASS. 단, present한 **선택필드(where/how/why) 중 `anchors=[]`인 것이 있음** — value는 있는데 앵커할 사실 토큰이 0개(구조적으로 검증 불가). 기계가 "이 필드는 미검증"이라 분류만. | state=REVIEW(잠금) |
| **FAIL**(시도 실패, 상태 아님) | 필수 누락, 또는 어떤 필드든 `anchors` 중 하나라도 원문에 없음(환각). | tries++ → TODO 유지(재시도). tries>Max면 state=DONE. |

> REVIEW는 **AI의 자기선언(interpretive=true)이 아니라 `anchors=[]`라는 구조 사실**로만 트리거된다.
> 즉 "해석이냐 사실이냐"를 주관 판정하지 않는다 — 앵커가 0개면 자동 REVIEW, 1개라도 있으면
> 그 토큰을 검증한다. `Field.interpretive`는 표시 라벨일 뿐 게이트 입력이 아니다.

## CLI 흐름
- `ccnews next` → 다음 TODO 기사 + (추출 본문 + 에이전트 프롬프트).
- 에이전트가 event6 제출 → 게이트가 원어 앵커 검증 → PASS/REVIEW/FAIL → 잠금/재시도.
- 종단(PASS/REVIEW) 도달 시 결과를 `--out` JSONL에 1회 append(Phase007).

## 열린 결정
- 앵커 매칭 정규화 수준(공백/조사/하이픈/유니코드 정규화 — 과탐/미탐 균형).
  단 앵커는 원문 표면형이므로 정규화는 **원문·앵커 양쪽에 동일 적용**(매핑 추론 아님).
- when 필드의 `value`는 ISO지만 게이트는 `anchors`의 원문 날짜 표면형만 본다(매핑 검증 안 함).
  연도 누락 표면형("6월 2일")의 value 연도 보정은 `extracted.published_at` 참고 — 산출 품질 사안이지 게이트 사안 아님.

## 다음 단계
Phase007(output) — 종단 기사를 JSONL로 증분 저장(status 필드로 collected/audit 구분).
