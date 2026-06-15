# reins

[![Version](https://img.shields.io/badge/version-v0.1.2-blue.svg)](https://github.com/park-jun-woo/reins/releases)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

**퀘스트 CLI 개발 프레임워크** (Go). 고삐(reins) — 완료 판정 권한을 AI에서 기계 게이트로 옮긴다.
*"생성은 확률적, 검증은 결정론적."* ([how-make-quest](https://www.parkjunwoo.com/tech/how-make-quest.md) 방법론의 재사용 구현)

AI 에이전트는 멀티스텝 작업에서 *자기 완료를 자기가 판정*해 조기 종료한다. reins는 완료 판정을
결정론적 **게이트**에 넘겨, 불안정한 생성기로도 신뢰 가능한 완료를 만든다. 에이전트는 일회용,
진행은 누적된다.

> **상세 API·패턴은 [MANUAL.md](MANUAL.md)** (Manual for AI Agents) 참조. 본 문서는 개요.

## 핵심 모델

- **래칫(ratchet)** — 한 방향 상태기계. 한 번 PASS면 불가역, 남은 일은 단조감소
  (`remaining(t+1) ≤ remaining(t)`).
- **게이트 = 규칙 카탈로그** — 위반 탐지 규칙의 집합. 각 규칙이 문제 발견 시 발동(true)하고 사실
  (`Fact`)을 싣는다. 심각도는 **레벨**(Fail/Review) — 가중치가 아니라 레벨이라 결정적 위반 1개가 곧
  FAIL. `Evaluate`가 발동 규칙을 레벨로 집계해 PASS/REVIEW/FAIL을 낸다.
- **권한 비대칭** — PASS 잠금은 **기계만**. L1 기계(결정론, PASS 단독권) / L2 AI(회의자, REVIEW만) /
  L3 사람(잔여).
- **사실 피드백** — FAIL은 의견이 아니라 위치·기대·실제값(`Fact`). 모델의 아첨 성향을 *수렴*으로
  돌린다.

## 두 게이트 백엔드

| 백엔드 | 언제 | 무엇 |
|---|---|---|
| **레벨집계** (`pkg/gate`) | 독립 규칙·단순 게이트 | `Rule` 카탈로그 + `Evaluate`. any Fail→FAIL, else Review→REVIEW, else PASS |
| **defeat 그래프** (`pkg/graph`) | 규칙 간 우선순위·근본원인 피드백 | toulmin h-Categoriser 백엔드. tautology PASS 워런트 + 위반 Counter + **`Supersedes`** 우선순위. 손으로 짠 가드를 선언적 엣지로 증발. 엣지0이면 레벨집계와 동치 |

defeat 그래프는 `Definition`이 `gate.Evaluator`를 구현하면 켜진다(opt-in). 부작용(HTTP/DNS)은
`pkg/ground` 원시연산 + **staged 평가**로 격리 — 싼 검사 실패 시 네트워크 fetch 자체가 안 일어난다.
그래프 평가는 에이전트 직통 **공략집**(`Verdict.Feedback`: "왜 졌나 + 뭘 바꿔야 이기나")을 낸다.

> toulmin은 `pkg/graph`·`pkg/ground`에만 결합되며 단방향(toulmin은 reins를 모름). `pkg/gate`·
> `pkg/cli`는 toulmin-free라, 그래프를 안 쓰는 소비자는 toulmin을 링크하지 않는다.

## 아키텍처 (`pkg/`)

| 패키지 | 역할 | 의존 |
|---|---|---|
| `pkg/textmatch` | 본문 포함 검증기 — `Normalize`(NFC)·`Contains`·`MissingTokens`. 환각 차단 원시연산 | x/text |
| `pkg/temporal` | 시간 명세 정규화 — 구조화 `Spec`(역법/성분/오프셋) → 그레고리력 ISO | (순수) |
| `pkg/quest` | 래칫 코어 — `State`·`Item`·`Verdict`/`Fact`·`Apply`·`Session`·`Export` | (순수) |
| `pkg/gate` | 게이트 계약 — `Definition`·`Rule`·`Context`·`Evaluate`(레벨집계)·`Evaluator`(그래프 훅) | quest |
| `pkg/graph` | defeat 그래프 백엔드 — `Graph`·`Warrant`·`Counter`·`Attacks`·`Supersedes`·`EvaluateStaged` | gate, quest, toulmin |
| `pkg/ground` | 네트워크 ground 원시연산 — `HTTPBody`·`MXResolves`(주입형 `Resolver`·요청당 스냅샷) | (순수 net) |
| `pkg/llm` | LLM 호출 어댑터 — ollama/xai/gemini chat completion. **생성(L0)만 담당, 판정/래칫과 무관**(권한 비대칭). num_ctx 자동 산정, 키는 env 전용 | net/http |
| `pkg/cli` | Cobra 스캐폴드 — `NewQuestCmd` → scan/next/submit/status/export/rules (+ 옵트인 `loop`) | cobra, quest, gate, llm |

## 명령 골격 (how-make-quest 정설)

```
scan    입력에서 N개 퀘스트 시드 + Progress 초기화 (스트리밍 소스는 소비자가 run 변형 추가)
next    TODO 하나 + 작성 프롬프트·검증 컨텍스트 출력
submit  제출 → 게이트 평가 → verdict → PASS 잠금 / FAIL이면 Fact·공략집 피드백
status  진행 집계 (PASS/REVIEW/DONE/TODO/SKIPPED/BLOCKED …)
export  종단 결과 JSONL 출력 (원본 보존, 1회 방출 래칫)
rules   게이트 규칙 카탈로그 출력 (자동 rulebook — 막는 치즈 목록 감사)
loop    (옵트인) submit의 자동 반복 — LLM이 생성, 게이트가 판정·잠금 (아래)
```

## 무인 자동 구동 — `loop` 명령 (옵트인)

외부 에이전트가 손으로 `next`→`submit` 하던 흐름을, **CLI 안에서 LLM이 생성하고 게이트가 판정하는 루프**로 닫는다. `cli.Options{Loop: &cli.LoopOptions{…}}`로 옵트인하면 `loop` 명령이 붙는다(미설정이면 미부착·완전 후방호환).

```
for 남은 TODO:
  system = 전역 + RuleSystem[직전 FAIL 근본원인 규칙]   # 규칙별 코칭
  raw    = backend.Complete(system, Render(it)+피드백)    # LLM 생성 (L0)
  verdict = 게이트 판정 → 래칫 Apply → export             # submit과 동일 경로
  FAIL이면 피드백(submit 출력과 동일) 되먹여 재시도(<MaxTries), 아니면 잠금→다음
```

- **권한 비대칭 유지** — LLM은 생성자(L0)일 뿐. **PASS 잠금은 여전히 게이트(기계)만**. MaxTries 초과 시 DONE 잠금 → 루프 단조 종료.
- **근본원인 규칙별 system 코칭** — `Verdict.RootCause`(평면·그래프 양쪽에서 결정론적 노출)로 직전에 못 넘은 규칙에 특화 지시를 되먹인다.
- **백엔드** — `--model ollama:gemma4:e4b`(기본) / `xai:…` / `gemini:…`. 로컬 ollama는 키 불요, num_ctx는 프롬프트 길이로 자동.

```bash
ccnews run --max-warcs 1                 # 시드(스트리밍 인제스천)
ccnews loop --model ollama:gemma4:e4b    # 남은 TODO를 gemma4가 생성→게이트 판정
```

## 퀘스트 만들기

`gate.Definition` 4개 메서드만 구현하면 reins가 래칫·명령 골격·집계·export를 공급한다:

```go
type Definition interface {
    Seed(args []string) ([]*quest.Item, error)            // 입력 → 초기 TODO 시드
    Render(s *quest.Session, it *quest.Item) (string, error)                // next가 보일 작성 프롬프트+검증 컨텍스트
    Prepare(s *quest.Session, it *quest.Item, raw []byte) (gate.Context, *quest.Verdict, error) // 제출 디코드 (short면 단락)
    Rules() []gate.Rule                                   // 게이트 위반-규칙 카탈로그
}

func main() { cli.NewQuestCmd("myquest", myDef{}, cli.Options{}).Execute() }
```

```go
// 치즈 방어 규칙 1개 = 위반 탐지기. 새 치즈 발견 → 규칙 하나 추가하면 게이트가 자란다.
var whoAnchorPresent = gate.Rule{
    Meta: gate.RuleMeta{ID: "who-anchor-present", Level: gate.LevelFail, Desc: "필수 who 앵커가 원문에 실재"},
    Check: func(ctx gate.Context) (bool, quest.Fact) {
        sub := ctx.Submission.(*Event6)
        if miss := textmatch.MissingTokens(ctx.Source, sub.Who.Anchors); len(miss) > 0 {
            return true, quest.Fact{Where: "who.anchors", Expected: "원문 substring", Actual: miss[0]}
        }
        return false, quest.Fact{}
    },
}
```

규칙 간 우선순위·근본원인 피드백·네트워크 검증이 필요하면 그래프 백엔드(`pkg/graph` + `gate.Evaluator`
+ `pkg/ground`)로 — [MANUAL.md](MANUAL.md) 참조.

## 상태

v1 빌드 완료 — 8패키지 `go build`·`go test` 통과, `filefunc validate` 0/0, `tsma` 0 TODO(전 함수
커버), gofmt clean. 레벨집계 + defeat 그래프 백엔드(toulmin) + ground/staged 평가 + **`loop` 무인
생성-검증 루프(`pkg/llm`, ollama/xai/gemini)** 까지 구현.

**첫 실사용 소비자 `comail`(이메일 수집)이 reins 그래프 게이트 위에서 종단 동작 확인** — scan/next/
submit/status/export, 그래프 판정·Supersedes·공략집 피드백·실네트워크 staged·래칫 잠금. (설계:
`plans/reins/Phase007-toulmin-gate.md`)

## 저장소 구성

- `pkg/` — 프레임워크 Go 모듈 (`github.com/park-jun-woo/reins`)
- `MANUAL.md` — AI 에이전트용 상세 매뉴얼
- `plans/reins/` — 설계문서 (Phase001~006 v1 스캐폴드, Phase007 toulmin 그래프 백엔드, Phase008 payload rehydration, Phase009 loop(구 agent) 명령·`pkg/llm`)
- `plans/ccnews/`·`plans/comail/` — 인스턴스 설계 + reins 이식 Phase
- `comail/`, `ccnews/` — 별도 모듈(자체 go.mod)인 퀘스트 인스턴스. reins를 import해 도메인만 구현
  (`comail`은 그래프 게이트로 이식 완료, `ccnews`는 이식 설계 중)

## 규약

- 결정론 게이트 명시 — 판정은 입력만으로, PASS 잠금은 기계만.
- 규칙은 위반 탐지기, 심각도는 레벨로 — 연속 가중은 그래프의 *진짜 경합*(L2 합의) 전용.
- 치즈 방어 우선 — "이 게이트를 어떻게 속이지?"의 답마다 규칙 하나(자동 카탈로그로 감사).
- 부작용은 ground로, 규칙은 순수 — 네트워크는 reins ground 원시연산이 소유, staged로 격리.
- N=1 추상화 금지 — 새 추상은 2번째 소비자(`ccnews`)로 검증 후 동결.
