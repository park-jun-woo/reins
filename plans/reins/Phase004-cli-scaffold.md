# Phase 004 — cli 스캐폴드: Cobra 명령 골격

- 프로젝트: reins (GitHub: `park-jun-woo/reins`)
- 작성일: 2026-06-05
- 상태: 스캐폴드 빌드됨 (pkg/cli)
- 의존: Phase001(헌장), Phase002(quest 코어), Phase003(gate 규칙 카탈로그)
- 패키지: `github.com/park-jun-woo/reins/pkg/cli`

## 목적

how-make-quest 정설 명령 골격 — **scan / next / submit / status / export** (+ reins 추가 `rules`)
— 을 Cobra로 한 번 구현한다. 새 퀘스트는 `gate.Definition`(Seed·Render·Prepare·Rules)만 끼우면
표준 CLI를 얻는다.

핵심 명제 한 줄: *"퀘스트마다 다른 건 규칙 카탈로그뿐. 명령 골격·상태기계·집계·출력은 reins가
구동한다."*

## 명령 골격

| 명령 | 역할 | 코어 호출 |
|---|---|---|
| `scan` | 입력에서 N개 퀘스트 시드 + Progress 초기화 | `def.Seed` → `Session.Save` |
| `next` | TODO 하나 + 작성 프롬프트·검증 컨텍스트 (읽기, 비변이) | `NextTODO` → `def.Render` |
| `submit` | 제출 → **규칙 카탈로그 평가 → 레벨 집계 verdict** → PASS 잠금 / FAIL이면 Fact 피드백 | `def.Prepare` → `gate.Evaluate` → `quest.Apply` → `Save` → `Export` |
| `status` | 진행 집계(TODO/PASS/REVIEW/DONE/SKIPPED/BLOCKED + TOTAL) | `Session.Progress` |
| `export` | 종단 결과 출력(원본 보존) | `quest.Export` |
| `rules` | 게이트 규칙 카탈로그 출력 (자동 rulebook — 치즈 방어 감사) | `gate.Catalog(def.Rules())` (RuleMeta) |

```go
package cli

// Options: Out(기본 "<name>-results.jsonl") + Version. 둘 다 선택.
type Options struct { Out string; Version string }

func NewQuestCmd(name string, def gate.Definition, opts Options) *cobra.Command
```

## submit 흐름 (규칙 카탈로그 집계)

```
세션 Find(key) → TODO 확인 → 제출물 raw 로드(--in 또는 stdin)
  → def.Prepare(it, raw) → (Context, short *Verdict)
  → short≠nil이면 그대로 verdict(SKIP 등), 아니면 gate.Evaluate(def.Rules(), ctx)
  → 레벨 집계: FAIL규칙↑→OutFail(Facts) / REVIEW규칙↑→OutReview / 무발동→OutPass
  → quest.Apply(item, verdict, now)            // 래칫 전이(now=UTC RFC3339)
  → Save → Export(terminal·미방출) → Save
  → "key -> OUTCOME (state STATE)" + FAIL이면 Facts 출력 (심볼릭 피드백)
```

## next↔submit 자기교정 루프

`submit` FAIL → 항목 TODO 유지 + **Fact 출력**(위치·기대·실제) → 다음 `next`가 같은 항목 재출력 →
에이전트가 사실 보고 수정 재제출. MaxTries 소진 시 `quest.Apply`가 DONE 잠금.

## scan vs 스트리밍 소스(run)

`scan`은 "입력에서 작업 시드"가 본질. 정적 입력(파일·CSV·디렉터리)은 1회 scan, 스트리밍·무한 소스
(ccnews CC-NEWS WARC)는 동일 시드를 인제스천 루프로 변형한 `run`(커서·래칫 재개). 둘 다 `def.Seed`
호출.

## 플래그 규약 (프레임워크 표준)

- **persistent**: `--session <path>` (기본 `session.json`), `--out <path>` (export; 기본
  `<name>-results.jsonl`).
- **submit 전용**: `--key <k>` (필수, 제출 대상 아이템), `--in <file>` (제출물 경로; `-` 또는 생략
  시 stdin).
- 도메인 플래그(`--cache-dir` 등)는 인스턴스가 추가.

## ccnews 매핑 (첫 통합)

| ccnews 현재 | reins cli |
|---|---|
| `cmd/run`(WARC 인제스천) | `scan`의 스트리밍 변형 `run` |
| `cmd/next` | `next` |
| `cmd/submit`(extract.Apply→anchor.Gate→Apply→Sweep) | `submit`(규칙 카탈로그 집계→Apply→Export) |
| (status 미구현) | `status` |
| `output.Sweep`/`--out` | `export` |
| (없음 — 게이트가 코드에 산개) | `rules` (앵커 게이트 카탈로그 자동 출력) |

## 결정론적 성격

cli는 IO·파싱·출력 포매팅만. 판정은 `gate`(결정론 규칙 집계), 전이는 `quest.Apply`(순수). cli 자체는
상태를 판정하지 않는다.

## 열린 결정

- `scan`/`run` 일반화: 하니스가 인제스천 루프 표준 제공 vs 스트리밍 시드는 인스턴스 몫(N=1 회피).
- `export` 포맷(JSONL 증분 기본 vs sink 인터페이스).
- `status`/`rules` 출력 형식(텍스트/JSON).
- 플래그 표준화 수준.

## 다음 단계

- 구현 완료: `pkg/cli` (Cobra scan/next/submit/status/export/rules) + `NewQuestCmd`·`Options`·JSONL
  sink (`quest_cmd_test.go`). quest/gate 연동.
- 통합: ccnews `cmd/*`를 `cli.NewQuestCmd`로 대체(도메인 Definition만 남김), 동작 보존 회귀.
- 2차 검증: comail을 같은 Definition으로 얹어 명령 골격·규칙 모델이 N=2에서 견디는지 확인 후 안정화.
