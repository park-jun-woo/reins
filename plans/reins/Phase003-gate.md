# Phase 003 — gate: 규칙 카탈로그 (핵심)

- 프로젝트: reins (GitHub: `park-jun-woo/reins`)
- 작성일: 2026-06-05
- 상태: 스캐폴드 빌드됨 (pkg/gate — 네이티브 레벨 집계; toulmin 비의존)
- 의존: Phase001(헌장), Phase002(quest 코어 — Verdict/Fact 소비)
- 패키지: `github.com/park-jun-woo/reins/pkg/gate`

## 목적

퀘스트를 퀘스트로 만드는 **결정론 검증기**를 **규칙 카탈로그**로 구현한다. 게이트는 단일
if-else가 아니라 **위반 탐지 규칙의 집합** — 각 규칙이 `//rule:` 어노테이션을 달고, 발동 시 사실
(Fact)을 싣고, 레벨로 심각도를 표한다. yongol이 동일 구조로 200+ 규칙을 운영해 입증한 패턴.

v1은 **toulmin 비의존 네이티브 레벨 집계**로 구현한다(`gate.Evaluate`). toulmin h-Categoriser는
*문서화된 미래 백엔드 플러그인 지점*이지 v1 의존이 아니다.

핵심 명제 한 줄: *"치즈 방어마다 규칙 하나. 규칙은 위반을 탐지하고, 레벨이 verdict를 가르고,
카탈로그가 게이트를 자기문서화한다 — '어떻게 속이지?'의 답이 grep되는 자산이 된다."*

## 규칙 모델 (네이티브 — pkg/gate)

규칙은 `Rule{Meta, Check}` 값이다. `Check`는 위반 탐지기 — `fired=true`면 문제 발견, evidence는
`quest.Fact`. `Context`는 제출 1건이 검사하는 사실을 운반한다.

```go
package gate  // rule.go

type Level int
const ( LevelFail Level = iota; LevelReview )   // String(): "FAIL"|"REVIEW"

type RuleMeta struct { ID string; Level Level; Desc string }   // 자동 rulebook 항목

// Context: 규칙이 검사하는 제출 사실. Source=캐시된 원천(치즈방어 재확인), Submission=디코드 산출물.
type Context struct { Item *quest.Item; Submission any; Source string }

// Rule: 위반 탐지기. Check가 fired=true + Fact를 내면 문제. 발동 규칙의 레벨 집계로 verdict 결정.
type Rule struct {
    Meta  RuleMeta
    Check func(ctx Context) (fired bool, fact quest.Fact)
}

// 예: 필수 앵커 원천 실재 (치즈 방어, Phase005 textmatch 재사용)
var requiredAnchorPresent = Rule{
    Meta: RuleMeta{ID: "who-anchor-present", Level: LevelFail, Desc: "필수 who 앵커가 원문에 실재"},
    Check: func(ctx Context) (bool, quest.Fact) {
        sub := ctx.Submission.(*Event6)
        if miss := textmatch.MissingTokens(ctx.Source, validAnchors(sub.Who.Anchors)); len(miss) > 0 {
            return true, quest.Fact{Where: "who.anchors", Expected: "원문 substring", Actual: miss[0]}
        }
        return false, quest.Fact{}   // 위반 없음
    },
}
```

- **발동(true) = 위반.** 아무 규칙도 안 터지면 PASS.
- **`//rule:` 어노테이션 + RuleMeta**: `ID`, `Level`(FAIL|REVIEW), `Desc`. 규칙 1개 = 파일 1개 →
  grep·codebook·자동 rulebook(`gate.Catalog`).
- **변주 = 클로저**: 임계값 변주는 같은 `Check`를 다른 파라미터로 닫는 클로저로(toulmin Spec 표면은
  백엔드 도입 시 노출 — v1 미구현).

## 레벨 집계 → verdict (yongol 방식, per-item)

게이트는 제출 1건에 대해 카탈로그를 평가하고 **발동한 규칙을 레벨로 집계**한다 — `gate.Evaluate`:

| 발동 상황 | verdict |
|---|---|
| FAIL-레벨 규칙 ≥1 발동 | **OutFail** (그 규칙들의 evidence = Facts) |
| FAIL 0, REVIEW-레벨 규칙 ≥1 발동 | **OutReview** |
| 아무 규칙도 안 터짐 | **OutPass** |
| (pre-gate) 신뢰/정책 단락 | **OutSkip / OutBlock** (`Definition.Prepare`의 short verdict) |

```go
package gate  // evaluate.go

// Evaluate: 카탈로그를 1건에 평가 → 발동 규칙을 레벨로 집계(가중치 아님). 발동 Fact는 규칙 ID로
// 스탬프되어 evidence로 수집(FAIL 시 에이전트에 환류). 결정론: 같은 (rules, ctx) → 같은 Verdict.
func Evaluate(rules []Rule, ctx Context) quest.Verdict

// Catalog: 규칙들의 RuleMeta — cli `rules`가 출력하는 자동 rulebook(막는 치즈 목록, ID·레벨).
func Catalog(rules []Rule) []RuleMeta
```

가중치 아님 — **레벨**. 단일 결정적 위반이 곧장 FAIL(가중 counter처럼 verdict 0=REVIEW로 새지
않음). h-Categoriser 연속 가중은 *진짜 경합*(아래 L2) 백엔드를 끼울 때만. SKIPPED/BLOCKED는
카탈로그 규칙이 아니라 `Definition.Prepare`가 반환하는 **short verdict**로 게이트를 단락한다.

## Definition — 퀘스트 도메인 계약 (4개 메서드)

인스턴스는 다음 인터페이스만 구현하면 reins가 래칫·명령 골격·집계·export를 공급한다:

```go
package gate  // rule.go

type Definition interface {
    Seed(args []string) ([]*quest.Item, error)           // 입력 → 초기 TODO 시드
    Render(it *quest.Item) (string, error)               // next가 보일 작성 프롬프트+검증 컨텍스트
    Prepare(it *quest.Item, raw []byte) (ctx Context, short *quest.Verdict, err error)
                                                          // 제출 디코드 → Context. short≠nil이면 게이트 단락(SKIP 등)
    Rules() []Rule                                        // 게이트 위반-규칙 카탈로그
}
```

## 게이트 5단 유도 = 규칙 패밀리 (how-make-quest)

게이트는 한 번에 안 짜고 다음 패밀리로 *유도*한다. 각 단계가 적절 레벨의 규칙을 보탠다:

| 단계 | 규칙 패밀리 | 레벨 |
|---|---|---|
| ① 형식 | 스키마·타입·필수필드 정합 | FAIL |
| ② 블랙리스트 | 명백한 쓰레기(`example.com`·플레이스홀더·빈값) | FAIL |
| ③ REVIEW 임계 | 모호 케이스 — 침묵 PASS 절대 금지 | REVIEW |
| ④ **치즈 방어** ★ | 핵심 사실 원천 재확인(환각·날조 차단) | FAIL |
| ⑤ 외부 일관성 | 네트워크·도메인(도달 실패는 REVIEW) | FAIL/REVIEW |

**④가 linchpin.** "이 게이트를 어떻게 속일까?"의 답마다 규칙 하나. 새 치즈 발견 = `//rule:` 파일
드롭(게이트가 자란다). 예: ccnews 빈앵커 `[""]`·플레이스홀더 value·환각 앵커·when 포맷 → 각 규칙.

## 예외 = defeat 엣지 (toulmin Except/Defeater — 설계 예약)

"규칙 발동, 단 X면 무효"는 toulmin 백엔드를 끼울 때 **공격 엣지**로 선언한다. yongol 실증:
`// @no-pagination`, `// @state-neutral`, `-- @func-managed` 가 해당 규칙을 defeat. 인스턴스가
정당한 예외를 어노테이션으로 표하면 위반이 취소된다 — if-else 중첩이 아니라 선언적 예외.

**v1 네이티브 집계에는 미구현.** v1에서 정당 예외는 `Check` 내부에서 직접 판정하거나(예외 조건을
보고 fired=false 반환) `Prepare`의 short verdict로 처리한다. defeat 그래프는 toulmin 도입 시 활성화.

## 권한 비대칭 — L1/L2/L3

| 계층 | 규칙 종류 | 권한 |
|---|---|---|
| **L1 기계** | 결정론 규칙(형식·블랙리스트·치즈방어·외부) | **PASS 단독권** |
| **L2 AI** | 독립 모델·좁은 예/아니오·복수 합의 | **REVIEW만** (h-Categoriser 가중 합의가 제값을 하는 유일 지점) |
| **L3 사람** | L1·L2 잔여 | 최종 |

같은 모델이 생성·판정하면 같은 맹점 → L2는 반드시 독립 회의자.

## 자동 rulebook (감사 자산)

RuleMeta(ID·Level·설명·source)에서 **게이트 카탈로그를 자동 생성**(yongol `rulebook.md`처럼). 모든
치즈 방어가 ID·레벨로 나열되어 "우리 게이트가 막는 치즈 목록"이 *문서가 아니라 코드에서 도출*된다.
cli `rules` 명령으로 출력(Phase004).

## 재사용 검증기

- **textmatch**(Phase005): ④ 치즈 방어 — 토큰이 원천에 글자 그대로 실재하는지.
- **temporal**(Phase006): ① 형식 + ④ 사실 — 시간 명세 정규화·성분 검증.

규칙 작성자는 이들을 *호출*한다(substring/날짜 로직 재구현 금지).

## 결정론적 성격

규칙은 (제출물 + 주입된 원천)만으로 발동을 결정하는 순수 함수(④의 원천 재확인은 캐시 입력 기준).
외부 IO는 ⑤로 분리, 도달 실패 시 REVIEW. 레벨 집계는 결정론. L2 AI만 REVIEW에 관여.

## 열린 결정

- `//rule:` 어노테이션 스키마(filefunc `//ff:`와 정합) + 자동 등록 방식.
- 레벨 enum: FAIL/REVIEW만 vs WARN(보고용) 추가.
- pre-gate(SKIPPED/BLOCKED) 규칙을 카탈로그에 포함 vs 별도 단계.
- L2 AI 회의자 호출을 프레임워크가 표준 지원할지 vs 인스턴스 몫.
- h-Categoriser 가중을 노출할 표면(경합/합의 케이스 전용).

## 다음 단계

- 구현 완료: `pkg/gate` — `Rule`·`Context`·`Definition`·`Evaluate`(레벨 집계)·`Catalog`·`RuleMeta`
  (`evaluate_test.go`). 네이티브(toulmin 비의존). defeat/h-Categoriser는 toulmin 백엔드 도입 시 추가.
- Phase004(cli) `submit`이 `Prepare`→(`Evaluate` 또는 short)→집계 verdict → `quest.Apply`.
- **첫 통합·검증**: ccnews 앵커 게이트(현 6~7조건: 빈앵커·플레이스홀더·환각·value위생·when포맷·
  SKIPPED·선택REVIEW)를 규칙 카탈로그로 시제품화 → 명령형 `check_field`/`check_required` 산개 대비
  작성·확장·감사 DX 비교. 우월하면 정식 채택, comail로 N=2 확인.
