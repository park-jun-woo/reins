# Phase 007 — toulmin 백엔드: defeat 그래프 게이트 + trace 피드백

- 프로젝트: reins (GitHub: `park-jun-woo/reins`)
- 작성일: 2026-06-05
- 상태: 설계 확정 대기

## 목적

게이트를 **평평한 규칙 카탈로그**에서 **defeat 그래프**로 끌어올린다. toulmin(논증 엔진,
h-Categoriser)을 reins의 판정 백엔드로 채택해 — **판정 · 근본원인 피드백 · 부작용 스케줄링**을
*한 모델*에서 떨어뜨린다.

핵심 명제: *"규칙은 독립이 아니다. 어떤 위반은 다른 위반을 무력화한다(defeat). 그 관계를 1급으로
선언하면, 게이트는 '무엇이 틀렸나'를 넘어 '무엇이 주범이고 무엇이 곁가지인가'를 안다 — 그리고
에이전트에게 PASS로 가는 공략집을 준다."*

## 배경 — comail이 드러낸 갭 (찍먹 1호)

첫 소비자 comail 이식([../comail/Phase001](../comail/Phase001-reins-adoption.md))에서 세 증상이 한 뿌리로 모였다:

1. **가드 = 손으로 짠 암묵적 defeat 그래프.** 규칙 6개에 `passedEarly`/freemail-no-fire/
   mismatch-no-fire 가드를 흩뿌렸다. 이는 "규칙 A가 발동하면 B·C는 무의미"를 *코드로 흉내낸* 것.
2. **평평한 Fact는 근본원인을 못 찾는다.** email-format도 틀리고 source-fetch도 실패하면, 게이트가
   엉뚱하게 source-fetch를 주범으로 보고 → 에이전트가 멀쩡한 출처를 붙잡고 헛고생.
3. **부작용 순서(G5).** `Prepare`가 형식검사보다 먼저 fetch → 형식 틀린 입력에도 네트워크 낭비.

셋 다 reins에 **규칙 우선순위/무력화를 표현할 수단이 없어서** 생긴다. defeat 그래프가 정확히 그
어휘다.

## 확정된 방향 (열린 결정 아님)

| 축 | 결정 |
|---|---|
| toulmin 범위 | **A — 풀 백엔드.** 판정+피드백 모두 toulmin 그래프로 (trace-only 레이어 아님) |
| 피드백 독자 | **무조건 에이전트.** 퀘스트는 에이전트가 처리한다. 사람/감사 분리 레인 없음 |
| 야심 수위 | **길게.** defeat 그래프로 *평가 스케줄링*까지 → G5(부작용 순서)도 같이 푼다 |
| 네트워크 검사 | **reins 기능으로 직접 제공** (ground provider 원시연산). 소비자가 fetch/MX를 손코딩하지 않는다. reins가 *소유*하므로 tier는 ground 사용에서 자동 파생 — 감지·딱지 불필요 |
| toulmin 결합 | **직접 import** (어댑터 없음). rego 등 대체재 부재 → 추상화는 YAGNI. 유일 제약: **단방향 의존**(toulmin은 reins를 모른다, 순환 금지) |

## 핵심 통찰 — 현 모델은 *엣지 0 그래프의 특수해*

이게 A를 "재작성"이 아니라 **상위집합 확장**으로 만든다.

- 현 reins(독립 규칙 + 레벨 집계) = **defeat 엣지가 0개인 그래프**.
- **책임 분리** (검수로 교정 — toulmin Strength에 의존하지 않는다):

| 레이어 | 누가 | 무엇 |
|---|---|---|
| **Level 메타** | **reins 보유** | 노드별 `Level`(Fail/Review) — toulmin 밖. 크리스프 판정의 근거 |
| **그래프 수학** | toulmin | defeat 엣지 + h-Categoriser verdict + **trace(노드별 Activated 기록)** |
| **Strength** | toulmin (v1 **미사용**) | Strict/Defeasible는 *경합 가중* 축. **공개 setter 부재**(`graph_rule.go`/`graph_counter.go`가 `Defeasible` 하드코딩, `Strict` 상수는 아무도 안 씀) + `TraceEntry`에 Strength 필드 없음 → v1은 toulmin Strength 미의존, reins Level로 판정 |

⟹ **엣지 없는 가게이트는 지금과 동일 동작**(단순 소비자 무피해). reins가 *trace의 Activated × 자체
Level맵*을 집계하면 현 `Evaluate`(any Fail→FAIL …)와 **동치**(엣지0·qualifier=1.0 전제에선
h-Categoriser float은 항상 +1 — 쓰이지 않음, 동치는 활성 판독으로 성립). supersession/엣지를 *추가*하면
근본원인·스케줄링이 켜진다. argumentation은 **opt-in**.

## 아키텍처

### ① 그래프형 `Definition`

`Rules() []Rule`(평평) → **그래프 선언**으로. **tautology PASS 워런트 1개 + 모든 위반 = 그 워런트를
공격하는 Counter** + 위반 간 우선순위 엣지.

> **토폴로지 주의(검수로 교정).** toulmin `isWarrant`는 *attacker 노드를 결과에서 제외*한다. 위반을
> Counter로만 두고 PASS 워런트를 안 만들면 결정타 위반이 `Evaluate` 결과에서 누락된다. 그래서 반드시
> **항상-활성 PASS 워런트**를 두고 위반들이 그걸 공격해야 한다. 그러면 `Evaluate`는 워런트 1개를
> 결과로 내고, 그 **trace에 카운터별 Activated가 실린다** → reins가 그걸 판독.

`gate.Graph`는 toulmin을 **직접 감싸는 reins 편의 빌더**(어댑터 추상 아님 — toulmin `g.Counter`/
`Attacks`를 그대로 호출하고, reins `Level`을 노드 메타로 부착한다. 실제 toulmin `*Rule` 빌더는
`Attacks`/`With`/`Qualifier` 3개뿐이라 `Level`은 reins 측 메타):

```go
func (d comailDef) Graph() *gate.Graph {
    g := gate.NewGraph("comail-email")
    pass := g.Warrant(alwaysTrue)                       // tautology: 항상 활성
    // 위반 = PASS 공격 카운터. 2번째 인자 = reins Level 메타(toulmin Strength 아님).
    fmtR   := g.Counter(ruleEmailFormat,       gate.Fail).Attacks(pass)
    holder := g.Counter(ruleSourceLacksEmail,  gate.Fail).Attacks(pass)
    mx     := g.Counter(ruleMxMissing,         gate.Fail).Attacks(pass)
    free   := g.Counter(ruleFreemail,          gate.Review).Attacks(pass)
    // 위반 간 우선순위 = reins-side **supersession**(toulmin 엣지 아님 — 아래 [중요] 참조)
    fmtR.Supersedes(holder, mx)                        // 형식 틀리면 출처·MX 카운터를 집계에서 제외
    free.Supersedes(holder)                            // 프리메일이면 출처검사 흡수 → REVIEW 보존
    return g
}
```

> **[중요] precedence ≠ toulmin defeat.** toulmin의 counter→counter `Attacks`는 verdict float만
>낮출 뿐 **`Activated`를 끄지 못한다**(`calc`: `active[id]`는 함수 결과로만 설정, defeat는 `sum`에만
> 반영). 그래서 위반 간 우선순위를 toulmin 엣지로 걸면 §② 판독(`Activated`)이 그걸 못 본다 — defeat된
> holder가 여전히 "활성 Fail"로 FAIL을 내고, freemail 케이스가 REVIEW→FAIL로 뒤집힌다. ⟹ **위반 간
> 우선순위는 reins-side `Supersedes` 메타**로 표현하고 §② 읽기 레이어가 크리스프하게 적용한다(상류
> 활성 카운터가 하류를 집계에서 제외). counter→pass `Attacks`만 toulmin 그래프(verdict/경합·Activated
> 기록·trace)에 남는다.

가드(`passedEarly` 등)가 **선언적 supersession으로 증발**한다 — A의 값을 N=1에서 즉시 검증할 지점.
(supersession은 판정·피드백·tier에 반영되나 *함수 실행*은 못 막음 — 부작용 차단은 §⑤ staged.)

### ② verdict 매핑 — 잠금 경계는 크리스프, 중간만 가중

h-Categoriser float `[-1,+1]`이 들어와도 **reins는 `trace.Activated` × 자체 Level맵으로 경계를
박는다**(float 임계값으로 PASS를 주지 않는다). 이유(엔진 수학 검증, `TestGraphWithDefeat`로 확인):
활성 카운터 1개(qual 1)가 워런트(qual 1)를 공격하면 워런트 verdict = `2·(1/(1+1))-1 = 0.0` — 강한
공격 1개는 워런트를 **0(미결)으로 끌 뿐 -1로 안 보낸다**. 그래서 "결정적 위반 1개 = FAIL"을
float으로는 못 만든다 → **카운터의 Activated와 reins Level을 직접 읽는다**(toulmin Strength 아님 —
공개 setter 부재).

판정 입력 = **활성 카운터에서 reins-superseded(상류 활성 카운터에 흡수된 것)를 뺀 잔존 카운터**. 그
잔존을 Level로 집계:

| 조건 (잔존 활성 카운터 × reins Level맵) | reins | 래칫 |
|---|---|---|
| 잔존 **Fail-Level** 카운터 ≥ 1 | **FAIL** (크리스프) | 재시도 / MaxTries→DONE |
| 잔존 카운터 0 (PASS 워런트 완승) | **PASS** | 불가역 잠금 |
| 잔존 카운터 있으나 **전부 Review-Level** | **REVIEW** | 사람 큐 |

→ **supersession이 곁가지를 크리스프하게 제외**(예: freemail이 holder를 superseded → 잔존=free만 →
REVIEW). 이건 reins-side·결정론이라 float 임계값 불요. 연속 가중 float은 *진짜 경합*(L2 합의,
qualifier를 0/1 밖으로 푼 노드 — comail엔 없음)에서만 REVIEW 밴드를 가른다. 즉 **precedence-only
게이트(comail)에선 h-Categoriser float이 휴면**하고 toulmin은 평가·trace 하니스로만 쓰인다 — float은
genuine contest 전용. "PASS 잠금은 기계만, 결정론" 보존. (검증: `calc`는 고정 qualifier·순수 ctx에
결정론·멱등 — `TestEvaluateIdempotent`)

### ③ trace → 에이전트 피드백 (PASS로 가는 공략집)

**설계 법칙(load-bearing)**: *trace의 잎(leaf)은 전부 결정론적 Fact. qualifier·specs는 고정
선언값.* 이게 지켜지면 trace는 논쟁거리가 아니라 **반박 불가능한 공략집**이 된다.

```
FAIL. 결정타 = email-format (잔존 활성 Fail 카운터, 상류).
  근거 Fact: where=email expected="user@domain 형식" actual="not-an-email"
  source-lacks-email은 email-format에 의해 superseded(상류 우선) → 집계·렌더에서 곁가지 처리.
  → 판정을 뒤집으려면 email-format을 끄라(유효 이메일 제출).
```

평평한 리스트("여기 틀림, 여기 틀림")가 아니라 **"왜 졌나 + 정확히 뭘 바꿔야 이기나"**가 한 구조에.
일회용 에이전트(수렴이 전부)에게 최적. (반논증 철학과 무모순: 의견이 아니라 *구조+사실*을 준다.)

### ④ 네트워크 = 프레임워크 기능 (ground provider 원시연산)

소비자가 fetch·MX를 손코딩하지 않는다. reins가 **검증 ground 원시연산**으로 직접 제공한다 —
`textmatch`/`temporal`이 글자·시간 원시연산이듯, **네트워크 ground**도 reins 자산:

| ground 원시연산 | 입력(소비자 공급) | 출력(reins 공급) |
|---|---|---|
| `HTTPBody(url)` | 출처 URL (제출에서) | 페이지 본문 (lazy·캐시) |
| `MXResolves(domain)` | 이메일 도메인 | 수신 가능 여부 bool |

- comail은 `fetch.go`·`lookup_mx.go`를 **버린다**. reins ground를 읽을 뿐.
- reins가 ground를 *소유*하므로 "어느 규칙이 네트워크를 쓰나"는 **선언으로 자명**(그 규칙이 어떤
  ground를 읽는지) — 감지·딱지·AST 마법 없음. tier 분리가 여기서 파생.
- ground는 **lazy**: 첫 읽힘에만 계산(부작용)·캐시. 단, "읽힘"은 규칙 함수 실행을 뜻하므로 lazy만으론
  부족 → 아래 ⑤ staged 평가가 *defeat된 규칙의 함수 자체*를 안 돌게 막는다.

### ⑤ 단계적(staged) 평가 → G5 통합

**검증된 엔진 사실** (`eval_context_calc.go` + `TestLazySkipsRebuttalWhenWarrantFalse`):
- toulmin 평가는 lazy — 워런트가 **비활성**이면 그 attacker는 calc 안 함(`if !ec.active[id] { return }`
  가드가 합산 전 반환). 즉 lazy가 무조건 다 돌리는 건 아니다.
- **그러나 우리 토폴로지의 PASS 워런트는 tautology(항상 활성)** → 그 **모든 카운터가 reached → 각
  함수 1회 실행**. 위반 간 supersession(reins-side)이나 toulmin float defeat는 집계·verdict만 바꿀
  뿐, holder가 pass의 attacker로 reached인 이상 **holder 함수 실행(=fetch)은 못 막는다**. ⟹ 단일
  평면 그래프에 무네트워크+네트워크 규칙을 섞으면 **네트워크 카운터의 fetch가 defeat돼도 발생.** G5는
  그래프 수학만으론 안 풀린다.

⟹ **reins가 단계 오케스트레이션을 소유한다**:

```
Tier 0 (무네트워크 그래프): email-format·placeholder·freemail·site-missing 평가
  └ terminal FAIL(활성 Fail-Level 카운터)이면 → 중단. Source 미접근 → fetch 안 일어남.
Tier 1 (네트워크 그래프): Tier0 생존 시에만. Source = lazy Ground(첫 접근 시 fetch).
  └ source-lacks-email·mx-missing·group-review 평가
```

각 tier는 toulmin 그래프 1회 평가. **부작용은 살아남은 tier에서만**. 근본원인 피드백·부작용 격리·
가드 제거가 **한 메커니즘**에서 떨어진다 — 야심 ③의 보상.

## 열린 결정

> 선결 2건은 **확정됨**(§확정된 방향): 네트워크 = reins ground provider 기능 / toulmin = 직접
> import·단방향. 아래는 잔여 미결.

1. **ground 의존 선언 형식** — 규칙이 어떤 ground(`HTTPBody`/`MXResolves`)를 읽는지 reins에 어떻게
   알리나? 명시 선언(`rule.Needs(HTTPBody)`) vs ground 접근 시점 기록. reins의 tier 위상정렬 입력.
2. **qualifier 고정 강제 범위** — 결정론 보존을 위해 qualifier=0/1 위주 강제. 연속 가중은 *언제*
   허용? (= L2 AI 합의 같은 진짜 경합에서만. 게이트가 "이건 경합 노드"라 선언한 곳.)
3. **그래프 acyclic 제약** — toulmin은 사이클 금지(`detectCycle`). 상호 무력화가 필요한 도메인이
   생기면 표현력 한계. 당장 comail엔 무관하나 기록.
4. **Definition API 브레이킹** — `Rules()` → `Graph()`. 엣지 없는 후방호환 경로(`Rules()`가
   엣지 0 그래프로 자동 승격)를 유지할지. 권장: 유지(상위집합 통찰의 실현).
5. **피드백 렌더 포맷** — trace→텍스트 변환 규약. 결정타 1개 + 무력화된 곁가지 표기 수위(완전 숨김 vs
   "defeat됨" 한 줄). 에이전트 토큰 예산과 수렴 속도의 트레이드오프.
6. **결정타 단일화** — 활성 Fail 카운터가 둘 이상이면 Fact를 다 줄까, defeat 그래프상 더 상류 1개만
   (진짜 근본원인) 줄까. 위상 최상류 선택 규칙 필요.
7. **toulmin 확장 vs reins-side 메타 [잠정: reins-side]** — toulmin은 Strict 공개 setter가 없고
   `TraceEntry`에 Strength가 없다. v1은 reins가 Level 메타를 자체 보유해 우회(채택). 장기엔 위반 간
   상호흡수에서 Strict가 필요해지면 toulmin에 `.Strict()`/`TraceEntry.Strength`를 추가 — **단방향
   의존 유지**(toulmin은 reins를 모름)이므로 toulmin 측 변경은 안전.
8. **네트워크 ground의 leaf 비결정성** — §③ "leaf=결정론 Fact" 법칙 vs `MXResolves`/`HTTPBody`의
   타임아웃·일시오류. ground는 **요청당 1회 스냅샷·캐시**해 그 평가 내 leaf를 고정(재시도는 다음 submit).
   스냅샷 실패(접근 불가)는 Tier 0 통과 후의 terminal FAIL Fact로 처리.
9. **precedence(supersession) vs genuine defeat(float) 경계** — 위반 간 *결정적 우선순위*는 reins-side
   `Supersedes`(크리스프, Activated 못 끄는 toulmin 엣지 대신). toulmin counter→counter `Attacks`(float
   defeat)는 *진짜 경합*(여러 Review 신호 가중)에만 쓴다. **comail은 전부 precedence라 float 휴면** —
   toulmin은 평가·trace 하니스로만. genuine-contest 소비자(L2)가 등장하면 float 경로 활성화. 두 메커니즘의
   API 표면(`.Supersedes` vs `.Attacks`)을 명확히 가를 것.

## 다음 단계 (Phase 분할 — 의존 체인)

| Phase | 범위 |
|---|---|
| 008 | **그래프 코어** — `gate.Graph`/`Warrant`/`Counter`/`.Attacks`(toulmin)/`.Supersedes`(reins-side). toulmin **직접 import**(단방향) + reins 편의 빌더. 노드별 **reins Level 메타 보유**(toulmin Strength 미의존). tautology PASS 워런트 패턴. 엣지 0 = 현 레벨집계 동치(후방호환 테스트) |
| 009 | **verdict 매핑** — `trace.Activated` × Level맵 → PASS/REVIEW/FAIL 크리스프 경계 + 래칫 배선. 결정론·멱등 보장 |
| 010 | **trace 피드백** — leaf=Fact 강제, 공략집 렌더러. 에이전트 직통 포맷 |
| 011 | **ground provider + staged 평가** — `HTTPBody`/`MXResolves` 원시연산, lazy 계산, tier 오케스트레이션. 부작용 격리(G5 해소) |
| 012 | **comail 재이식 (N=1 증명대)** — 8규칙+가드 → 엣지 5~6개 그래프, `fetch.go`/`lookup_mx.go` 삭제(reins ground 사용). 가드 증발·근본원인 피드백·fetch 생략을 실측. A의 가치 검증 후 안정화 |

원칙: 008~011은 reins 코어, 012가 *증명*. comail이 그래프 모델을 견디지 못하면 추상화를 후퇴시킨다
(N=1 추상화 금지 — 2번째 소비자까지 본 뒤 동결).
