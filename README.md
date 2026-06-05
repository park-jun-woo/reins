# reins

**퀘스트 CLI 개발 프레임워크** (Go). 고삐(reins) — 완료 판정 권한을 AI에서 기계 게이트로 옮긴다.
*"생성은 확률적, 검증은 결정론적."* ([how-make-quest](https://www.parkjunwoo.com/tech/how-make-quest.md) 방법론의 재사용 구현)

AI 에이전트는 멀티스텝 작업에서 *자기 완료를 자기가 판정*해 조기 종료한다. reins는 완료 판정을
결정론적 **게이트**에 넘겨, 불안정한 생성기로도 신뢰 가능한 완료를 만든다. 에이전트는 일회용,
진행은 누적된다.

## 핵심 모델

- **래칫(ratchet)** — 한 방향 상태기계. 한 번 PASS면 불가역, 남은 일은 단조감소
  (`remaining(t+1) ≤ remaining(t)`).
- **게이트 = 규칙 카탈로그** — 위반 탐지 규칙의 집합(yongol/toulmin 패턴). 각 규칙이 문제 발견 시
  발동(true)하고 사실(`Fact`)을 싣는다. 심각도는 **레벨**(Fail/Review) — 가중치가 아니라 레벨이라
  결정적 위반 1개가 곧 FAIL이다. `Evaluate`가 발동 규칙을 레벨로 집계해 PASS/REVIEW/FAIL을 낸다.
- **권한 비대칭** — PASS 잠금은 **기계만**. L1 기계(결정론, PASS 단독권) / L2 AI(회의자, REVIEW만) /
  L3 사람(잔여).
- **사실 피드백** — FAIL은 의견이 아니라 위치·기대·실제값(`Fact`). 모델의 아첨 성향을 *수렴*으로
  돌린다.

## 아키텍처 (`pkg/`)

| 패키지 | 역할 | 의존 |
|---|---|---|
| `pkg/textmatch` | 본문 포함 검증기 — `Normalize`(NFC+공백)·`Contains`·`MissingTokens`. 환각 차단 원시연산 | x/text |
| `pkg/temporal` | 시간 명세 정규화 — 구조화 `Spec`(역법/성분/오프셋) → 그레고리력 ISO. 미정은 `Determined=false` | (순수) |
| `pkg/quest` | 래칫 코어 — `State`·`Item`·`Verdict`/`Fact`·`Apply`·`Session`·`Export` | (순수) |
| `pkg/gate` | 규칙 카탈로그 — `Rule`·`RuleMeta`·`Context`·`Definition`·`Evaluate`(레벨 집계)·`Catalog` | quest |
| `pkg/cli` | Cobra 스캐폴드 — `NewQuestCmd` → scan/next/submit/status/export/rules | cobra, quest, gate |

검증기(textmatch/temporal)·quest 코어는 외부 비의존 순수. toulmin h-Categoriser(가중 논증)·defeat
엣지는 *진짜 경합/L2 합의*용 **미래 백엔드 플러그인 지점**으로 예약(v1 미의존).

## 명령 골격 (how-make-quest 정설)

```
scan    입력에서 N개 퀘스트 시드 + Progress 초기화 (스트리밍 소스는 run 변형)
next    TODO 하나 + 작성 프롬프트·검증 컨텍스트 출력
submit  제출 → 규칙 카탈로그 평가 → 레벨 집계 verdict → PASS 잠금 / FAIL이면 Fact 피드백
status  진행 집계 (PASS/REVIEW/DONE/TODO …)
export  종단 결과 JSONL 출력 (원본 보존, 1회 방출 래칫)
rules   게이트 규칙 카탈로그 출력 (자동 rulebook — 막는 치즈 목록 감사)
```

## 퀘스트 만들기

`gate.Definition` 4개 메서드만 구현하면 reins가 래칫·명령 골격·집계·export를 공급한다:

```go
type Definition interface {
    Seed(args []string) ([]*quest.Item, error)            // 입력 → 초기 TODO 시드
    Render(it *quest.Item) (string, error)                // next가 보일 작성 프롬프트+검증 컨텍스트
    Prepare(it *quest.Item, raw []byte) (gate.Context, *quest.Verdict, error) // 제출 디코드 (short verdict면 단락)
    Rules() []gate.Rule                                   // 게이트 위반-규칙 카탈로그
}

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

func main() { cli.NewQuestCmd("myquest", myDef{}, cli.Options{}).Execute() }
```

## 상태

스캐폴드 v1 빌드 완료 — `go build ./...` · `go test ./...` 통과, `filefunc validate` 0 error/0 warn,
tsma 0 TODO(전 함수 PASS/DONE), gofmt clean.

## 저장소 구성

- `pkg/` — 프레임워크 Go 모듈 (`github.com/park-jun-woo/reins`)
- `plans/reins/` — 설계문서 (Phase001 헌장 + quest-core·gate·cli·textmatch·temporal)
- `plans/ccnews/` — 인스턴스 ccnews 설계문서
- `ccnews/`, `comail/` — 별도 모듈(자체 go.mod)인 퀘스트 인스턴스. reins를 import해 도메인 규칙만 구현
  (`.ffignore`/`.tsmignore`로 reins 검증 대상에서 제외)

## 규약

- 결정론 게이트 명시 — 판정은 입력만으로, PASS 잠금은 기계만.
- 규칙은 위반 탐지기, 심각도는 레벨로 — 가중치로 하드체크를 흉내내지 않는다.
- 치즈 방어 우선 — "이 게이트를 어떻게 속이지?"의 답마다 규칙 하나(자동 카탈로그로 감사).
- N=1 추상화 금지 — 하니스/CLI는 ccnews 증류 v1, comail(2번째 소비자)로 검증 후 안정화.
