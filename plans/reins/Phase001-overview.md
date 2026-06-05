# Phase 001 — reins 개요 & 프레임워크 헌장

- 프로젝트: reins (GitHub: `park-jun-woo/reins`)
- 작성일: 2026-06-05
- 상태: 스캐폴드 빌드됨 (pkg/ — textmatch·quest·gate·temporal·cli)

## 목적

reins는 **퀘스트 CLI 개발 프레임워크**다(Go + Cobra). `how-make-quest`의 방법론 —
*"생성은 확률적, 검증은 결정론적"* — 을 재사용 코드+규약으로 구현해, 매번 골격을 다시 짜지 않고
신뢰 가능한 완료를 보장하는 퀘스트 CLI를 짓게 한다.

핵심 명제 한 줄: *"고삐(reins): 완료 판정 권한을 AI에서 기계 게이트로 옮긴다. 게이트는
**규칙 카탈로그** — 위반 탐지 규칙의 레벨 집계로 PASS/REVIEW/FAIL을 결정론적으로 낸다."*

## 퀘스트의 5대 불가결 요소 (how-make-quest)

reins가 ②~⑤의 *기계*를 제공하고, 인스턴스는 ①과 도메인 규칙만 짠다.

| # | 요소 | 누가 제공 |
|---|---|---|
| ① | **Goal** | 인스턴스 |
| ② | **Completion Condition** | reins(상태기계) + 인스턴스 **규칙 카탈로그** |
| ③ | **Verifier / Gate** | reins(규칙 엔진·레벨 집계·검증기) + 인스턴스 규칙 |
| ④ | **Feedback** | reins(Fact = 규칙 evidence) |
| ⑤ | **Progress** | reins(래칫·세션) |

## 두 기둥

**래칫 (한 방향 상태기계)** — 한 번 PASS면 불가역, 남은 일 단조감소(`remaining(t+1) ≤
remaining(t)`). **에이전트 일회용, 진행 누적.**

**권한 비대칭 (L1/L2/L3)** — PASS 잠금은 **기계만**. L1 기계(결정론, PASS 단독권) / L2 AI(회의자,
REVIEW만) / L3 사람(잔여).

## 게이트 모델 — 규칙 카탈로그 (yongol 입증)

reins 게이트는 단일 if-else 함수가 아니라 **위반 탐지 규칙의 카탈로그**다. yongol이 동일 구조로
200+ 규칙을 운영해 입증했다. v1은 toulmin 비의존 **네이티브 레벨 집계**로 구현(`gate.Evaluate`);
toulmin h-Categoriser는 문서화된 *미래 백엔드 플러그인 지점*이다.

- **규칙 = 위반 탐지기**: `Rule{Meta RuleMeta; Check func(Context) (fired bool, fact quest.Fact)}`.
  발동(true) = 문제 발견.
- **`//rule:` 어노테이션 + RuleMeta**(ID·Level·설명) → 규칙 1개 = 파일 1개, **카탈로그 자동 문서화**.
- **레벨 집계로 verdict**: FAIL-레벨 규칙 발동 → FAIL(evidence=Fact); 아니면 REVIEW-레벨 발동 →
  REVIEW; 아무것도 안 터지면 → PASS. (yongol ERROR/WARNING 집계를 *제출 단위*로 적용)
- **예외 = defeat 엣지(설계 예약)**: "규칙 발동, 단 `@override` 있으면 무효" — toulmin Except/Defeater
  백엔드를 끼울 때 활성화. v1 네이티브 집계에는 미구현.
- **치즈 방어 = 규칙 추가**: 새 치즈벡터 발견 → `//rule:` 파일 하나 드롭. 게이트가 *자란다*.

가중 논증(h-Categoriser)은 *진짜 경합·합의*(L2 AI 복수 투표) 같은 비결정 케이스에만 — toulmin 백엔드를
끼울 때. 하드 체크는 가중치가 아니라 **레벨**로 표현한다(단일 결정요인이 verdict 0=REVIEW로 새는 함정 회피).

## 명령 골격 (정설)

```
scan    입력에서 N개 퀘스트 시드 + Progress 초기화 (스트리밍 소스는 run 변형)
next    TODO 하나 + 작성 프롬프트·검증 컨텍스트 출력 (읽기)
submit  제출 → 규칙 카탈로그 평가 → 레벨 집계 verdict → PASS 잠금 / FAIL이면 Fact 피드백
status  진행 집계
export  종단 결과 출력 (원본 보존)
rules   게이트 규칙 카탈로그 출력 (자동 rulebook — 치즈 방어 감사)
```

next↔submit은 자기교정 쌍: submit FAIL → TODO 유지 + Fact → next가 다시 줌 → 수정 재제출.

## 아키텍처 계층 (의존: 위→아래)

```
github.com/park-jun-woo/reins  (Go module; 패키지는 pkg/ 아래)
├─ pkg/textmatch  순수 검증기 — Normalize(NFC+공백)+Contains+MissingTokens.   (Phase005)
├─ pkg/temporal   순수 검증기 — 시간 명세 → 그레고리력 ISO 정규화.            (Phase006)
├─ pkg/quest      래칫 코어 — 상태기계·Session·Progress·Verdict/Fact·          (Phase002)
│                 Apply·Export. Cobra·toulmin·도메인 비의존.
├─ pkg/gate       규칙 카탈로그 — RuleMeta·Rule·Context·Definition·             (Phase003)
│                 Evaluate(레벨 집계)·Catalog. textmatch/temporal을 규칙이 사용.
│                 네이티브 구현(toulmin 비의존; toulmin은 미래 백엔드 플러그인).
└─ pkg/cli        Cobra 스캐폴드 — scan/next/submit/status/export/rules.        (Phase004)
```

import 경로는 `github.com/park-jun-woo/reins/pkg/<패키지>` (예: `.../pkg/quest`). 인스턴스
(ccnews·comail)는 reins를 import해 **Goal + `gate.Definition`(`//rule:` 규칙 카탈로그)만** 짜면 완성.

## 정체성

- **이다**: Go 라이브러리(래칫·규칙엔진·검증기) + 규약(문서표준·결정론 게이트·`//rule:` 어노테이션).
- **아니다**: 런타임·코드젠·DSL 같은 무거운 프레임워크.
- **의존**: cli=Cobra(spf13), textmatch=x/text(NFC). quest 코어·gate·temporal은 외부 비의존(순수
  표준 라이브러리). toulmin은 v1에서 의존하지 않는다(미래 가중 백엔드 플러그인 지점으로만 예약).

## Phase 인덱스

| Phase | 문서 | 범위 |
|---|---|---|
| 001 | `overview` | 본 문서 — 프레임워크 헌장 |
| 002 | `quest-core` | 래칫 상태기계·Session·Progress·Verdict/Fact·Apply·Export (순수 코어) |
| 003 | `gate` | 규칙 카탈로그 — RuleMeta·Rule·Context·Definition·레벨 집계·자동 rulebook |
| 004 | `cli-scaffold` | Cobra scan/next/submit/status/export/rules |
| 005 | `textmatch` | 정규화(공백+NFC)+substring 포함 검증 (치즈 방어 규칙 헬퍼) |
| 006 | `temporal` | 구조화 시간 명세 → 그레고리력 ISO 정규화 (규칙 헬퍼) |

## 설계 원칙 (헌장)

1. **결정론 게이트** — 판정은 입력만으로. PASS 잠금은 기계만.
2. **규칙은 위반 탐지기, 심각도는 레벨로** — 가중치로 하드체크를 흉내내지 않는다.
3. **치즈 방어 우선** — "어떻게 속이지?"의 답마다 규칙 하나(자동 카탈로그로 감사).
4. **다국어 읽기는 AI, 변환·비교는 기계** — 라이브러리는 자유텍스트 파싱을 안 떠안는다.
5. **N=1 추상화 금지** — 하니스/CLI는 ccnews 증류 v1. comail(2번째 소비자)로 검증 후 안정화.

## 열린 결정 / 다음 단계

- gate는 v1에서 네이티브 레벨 집계로 구현(toulmin 비의존) — 확정. toulmin h-Categoriser는 진짜
  경합·합의(L2) 백엔드로 끼울 때만 도입 — 열림.
- 스캐폴드(pkg/) 빌드 완료: textmatch·quest·gate·temporal·cli + 단위테스트.
- 첫 통합: ccnews 앵커 게이트를 `gate.Definition` 규칙 카탈로그로 재구성. 2차 검증: comail로 N=2.
