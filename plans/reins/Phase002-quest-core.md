# Phase 002 — quest 코어: 래칫·세션·진행

- 프로젝트: reins (GitHub: `park-jun-woo/reins`)
- 작성일: 2026-06-05
- 상태: 스캐폴드 빌드됨 (pkg/quest)
- 의존: Phase001(헌장)
- 패키지: `github.com/park-jun-woo/reins/pkg/quest`

## 목적

퀘스트의 **불가역 진행 기계** — 한 방향 상태기계(래칫), 세션 영속, 진행 추적, verdict 적용, 사실
(Fact) 운반, 종단 export — 를 도메인·toulmin·Cobra 비결합 순수 코어로 제공한다. how-make-quest의
5요소 중 **Completion Condition·Progress·Feedback·verdict 적용**의 기계 부분.

핵심 명제 한 줄: *"에이전트는 일회용, 진행은 누적. 한 번 PASS면 다시 안 풀린다 —
`remaining(t+1) ≤ remaining(t)`."*

## 상태기계 (래칫 — 단방향)

```
TODO ──PASS──►  PASS     (게이트 통과 — 불가역 잠금)
     ──REVIEW─► REVIEW   (모호 — 사람 판단, terminal)
     ──FAIL───► tries++ ; tries≥MaxTries → DONE   (소진 종료)
     ──SKIP───► SKIPPED  (신뢰 불가 — 게이트 미실행, terminal)
     ──BLOCK──► BLOCKED  (정책 거부, terminal)
```

terminal 도달 시 재오픈 없음. `remaining = count(TODO)` 단조감소 → 진행 비진동 보장.

## 타입 (실제 코드 — pkg/quest)

`State`(아이템 래칫 위치)와 `Outcome`(게이트 1회 판정)을 **구분**한다: FAIL은 재시도 가능한 실패
시도이지 아이템 종단 상태가 아니다.

```go
package quest

// State: 아이템 래칫 위치. terminal 상태는 잠금(state.go).
type State string // TODO|PASS|REVIEW|DONE|SKIPPED|BLOCKED
func (s State) Terminal() bool   // PASS/REVIEW/DONE/SKIPPED/BLOCKED

// Outcome: 게이트의 제출 1건 판정. State와 별개(verdict.go).
type Outcome string
const ( OutPass Outcome = "PASS"; OutReview = "REVIEW"; OutFail = "FAIL"
        OutSkip = "SKIPPED"; OutBlock = "BLOCKED" )

// Fact: 게이트 규칙이 발동 시 싣는 정량·위치지정 evidence (how-make-quest "아첨할 틈 없는 사실").
type Fact struct { Rule, Where, Expected, Actual string }   // item.go

// Verdict: 게이트(Phase003)의 규칙 레벨 집계 결과.
type Verdict struct {
    Outcome Outcome   // OutPass|OutReview|OutFail|OutSkip|OutBlock
    Facts   []Fact    // FAIL/REVIEW를 일으킨 규칙들의 evidence
}
func (v Verdict) Reason() string   // Facts를 감사 로그용 한 줄로 렌더(verdict.go)

type Attempt struct { Try int; Outcome string; Reason string }   // 감사 로그 1줄

type Item struct {
    Key         string
    State       State
    Tries       int
    Verdict     string     // 마지막 Outcome 문자열
    Reason      string     // 마지막 Verdict.Reason()
    CollectedAt string     // 잠금 시각(now 주입)
    Log         []Attempt
    Emitted     bool       // export 1회 보장 래칫
    Payload     any        // 도메인 산출물(event6 등)
}

type Session struct { Version int; Items []*Item }   // session.go
func New() *Session                                  // Version=1 빈 세션
func Load(path string) (*Session, error)             // 부재는 os.IsNotExist(err)
func (s *Session) Save(path string) error            // pretty JSON
func (s *Session) NextTODO() *Item
func (s *Session) Find(key string) (*Item, error)
func (s *Session) Progress() map[State]int

const MaxTries = 3

// Apply: verdict를 상태기계에 적용. PASS/REVIEW/SKIPPED/BLOCKED → 잠금(+CollectedAt),
// FAIL → Tries++ (≥MaxTries면 DONE 잠금). Attempt 로그 항상. now 주입(순수, 판정 로직 없음).
func Apply(it *Item, v Verdict, now string)

// Sink: export 수신 인터페이스(포맷은 구현이 선택; cli가 JSONL sink 제공).
type Sink interface { Emit(it *Item) error }

// Export: terminal·미방출 항목을 sink로(원본 보존) + Emitted 래칫. 증분, 1회 보장.
func Export(s *Session, sink Sink) (n int, err error)
```

## 코어의 책임 경계

- **판정 안 함**: 코어는 전이·카운트·잠금·export만. PASS/REVIEW/FAIL *결정*은 게이트(Phase003)의
  규칙 카탈로그가, 그 verdict를 *적용*만 코어가.
- **Fact 운반**: 게이트 규칙이 만든 `[]Fact`를 verdict에 실어 받아, FAIL 시 cli가 그대로 에이전트에
  피드백(심볼릭 수렴 루프). 코어는 Fact를 해석하지 않고 운반·기록.
- **순수**: `Apply`/집계 로직은 입력만으로 결정. IO(Load/Save/Export)는 얇게 분리.

## 결정론적 성격

`Apply`는 (item, verdict, now)만으로 결정되는 순수 상태 변이. `Export`는 terminal·미방출만 1회
방출(Emitted 래칫). 래칫 단방향성·MaxTries→DONE·1회 방출이 핵심 불변식.

## 열린 결정

- `Payload any` vs 제네릭 `Item[T]`.
- `MaxTries` 상수 vs 퀘스트별 주입.
- `Sink` 표면(JSONL 증분 기본 vs 인스턴스 포맷 선택).
- SKIPPED/BLOCKED를 게이트 전 단계(pre-gate 규칙)로 둘지 코어 상태로 직접 둘지(Phase003과 연동).

## 다음 단계

- 구현 완료: `pkg/quest` + 단위테스트(`apply_test.go` — 래칫 단방향성, MaxTries→DONE, Export 1회 보장).
- Phase003(gate)이 Verdict/Fact 생산 → 코어 Apply가 소비.
- ccnews `internal/session`+`anchor.Apply`+`output.Sweep`를 이 코어로 흡수(첫 통합).
