# Phase 010 — when 필수 해제 (발행일은 구조화 데이터로 직공급)

- 프로젝트: ccnews (GitHub: `park-jun-woo/quest-ccnews`)
- 작성일: 2026-06-04
- 상태: 설계 — 구현 대기
- 의존: Phase005(구조화 추출), Phase006(event6 앵커 게이트), Phase009(앵커 강화)

## 목적

event6의 `when`을 **필수(who/when/what)에서 선택으로 내린다.** 발행일은 AI가 본문에서
추출·앵커할 대상이 아니라 **구조화 데이터(`datePublished`)에서 기계가 결정론적으로 얻는 값**이기
때문이다. AI에게 본문 앵커를 강요하면 날짜가 본문에 없을 때 헛앵커만 난다.

핵심 명제 한 줄: *"발행일은 추출 문제가 아니라 이미 가진 메타데이터다. when을 AI 필수 앵커에서
빼고, 날짜는 구조화 데이터로 직공급한다."*

## 배경 (Haiku/Sonnet 실측, Phase009 강화 게이트 위)

강화 게이트로 100건 샘플을 Haiku·Sonnet으로 처리한 결과, 빈앵커·플레이스홀더 치즈는 0이 됐으나
**`when` 앵커-값 단절**이 남았다. 날짜 값인데 앵커가 날짜를 지지 못하는 케이스:

| 모델 | when 의심(날짜값·숫자앵커 없음) | 명백한 가짜 앵커 |
|---|---|---|
| Haiku | 6/24 PASS | `when="2026-06-03" anchors=["Türkiye"]` 등 |
| Sonnet | 4/16 PASS | `when="2026-06-04" anchors=["Kurtulmuş"]`, `when="2026" anchors=["वायरल"]` |

**강한 모델도 동일 양상** — 게이트가 "앵커 존재"만 보고 "값 지지"는 못 보는 사각지대라
모델 불변(model-invariant)이다. 명시적으로 "날짜를 인명/지명에 걸지 말라"고 지시한 Sonnet도
본문에 날짜 토큰이 없으면 헛앵커를 냈다. 원인은 모델 역량이 아니라 **when을 필수로 둔 설계**다.

## 진단 — 왜 유독 when만 무너지나 (3 구조적 원인)

| # | 원인 | 설명 |
|---|---|---|
| 1 | **날짜가 본문이 아닌 메타에 산다** | 발행일은 JSON-LD `datePublished`·`<meta>`에 있고 본문 산문엔 없을 때가 잦다. 앵커 게이트는 본문 텍스트만 대조하므로 매칭 대상에 날짜 토큰이 없다. |
| 2 | **상대 표현 + 정규화** | `오늘`·`수요일`·`전날` → value는 절대 ISO(`2026-06-04`). 정직한 앵커(`오늘`)엔 그 숫자가 없어 verbatim substring이 안 잇는다. 6필드 중 정규화가 가장 심한 필드. |
| 3 | **필수 + 앵커 강제** | who/what은 본문에 글자로 박혀 문제없으나, when은 날짜 토큰이 없어도 `≥1 유효앵커`를 내야 한다 → 게이트 요건이 모델을 헛앵커로 떠민다. |

## 코드 실태 (구현 전 반드시 인지)

발행일 매핑 경로는 **이미 존재**한다:

```
JSON-LD datePublished | OG article:published_time
   → extract.Fields.PublishedAt            (map_article.go:11 / apply_meta.go:17-18)
   → session.Extracted.PublishedAt         (extract.Apply, apply.go:30-36, PublishedAt=L33)
   → output.Record.published_at            (render_collected.go:19-20)
```

**그러나 `extract.Apply`(= `a.Extracted`를 채우고 구조화 데이터 없으면 SKIPPED를 거는 함수)는
파이프라인에 미배선이다.** `cmd/submit.go:70`·`cmd/next.go:46`은 `extract.Parse`(앵커 대상
본문만)만 호출하고, `a.Extracted = …` 대입은 **`extract.Apply`(apply.go:30)가 유일한 출처인데
이 함수는 `*_test.go`에서만 불린다**(`cmd/`·`internal/ingest/`에서 호출 0건). 주의: `cmd/submit.go:73`의
`anchor.Apply`는 게이트 판정을 상태기계에 적용하는 **별개 함수**이지 이 추출기가 아니다. 결과:

- 현재 `a.Extracted`는 항상 nil → `render_collected.go:19`의 `if a.Extracted != nil`이 늘 거짓 →
  출력 `published_at`이 **항상 빈칸**.
- Phase005 SKIPPED 게이트(`extract.Gate`: 구조화 출처/제목 없음 → SkipNoStructured)도 **휴면** —
  비구조화 기사도 SKIPPED 없이 그대로 앵커 게이트행.

즉 "발행일 기계 직공급" 전제는 추출 코드 차원에선 참이나 **배선이 빠져 출력엔 안 실린다.**
따라서 when 제거는 이 배선(B)과 **반드시 함께** 가야 날짜를 잃지 않는다.

## 변경 (결정론적)

### A. when 필수 해제 — `internal/anchor/gate.go`

```go
// 변경 전
required := []namedField{ {"who", ev.Who}, {"when", ev.When}, {"what", ev.What} }
optional := []namedField{ {"where", ev.Where}, {"how", ev.How}, {"why", ev.Why} }

// 변경 후
required := []namedField{ {"who", ev.Who}, {"what", ev.What} }
optional := []namedField{ {"when", ev.When}, {"where", ev.Where}, {"how", ev.How}, {"why", ev.Why} }
```

게이트 의미(나머지 로직 불변, Phase006·009 그대로):
- **필수 who/what**: 값 위생(Phase009 L3) 통과 + 유효앵커(Phase009 L0) ≥1 + 전부 substring. 누락·환각 → FAIL.
- **선택 when/where/how/why**: nil이면 무시. present면 값 위생+앵커 검사 — 환각 → FAIL, 유효앵커 0 → **REVIEW**.
- when에 날짜 토큰이 없으면 AI는 **그냥 생략(nil)** 하면 됨 → 헛앵커 강제 소멸. 발행일은 아래 B가 별도 보유.

### B. 발행일 직공급 배선 — `extract.Apply` 연동 (필수 동반)

- `cmd/submit`(및 정석적으로 `run` 시드/`next`)에서 `extract.Apply(a, htmlBytes)`를 호출해
  `a.Extracted{Source, PublishedAt, …}`를 채우고, 구조화 데이터가 없으면 **SKIPPED로 잠근다**
  (Phase005 의도 활성화). `extract.Apply`는 PASS 시 앵커 대상 본문 `bodyText`를 반환하므로
  (apply.go:40) 현재 submit의 `extract.Parse`→`res.BodyText` 호출(submit.go:70·72)을 흡수·대체한다.
  단 **`ok=false`(SKIPPED) 분기에서는 submit이 `anchor.Gate`를 건너뛰고 즉시 저장·반환**하도록
  제어흐름을 함께 손봐야 한다(이미 잠긴 기사에 빈 본문으로 게이트를 돌리지 않게). 단순 함수 치환이
  아니라 SKIPPED 단락(short-circuit) 분기 추가가 포함된다.
- 효과: 출력 `published_at`이 구조화 출처(신뢰 100%)로 채워지고, 레코드는 **AI `when` 없이도
  항상 날짜를 가진다.** when 제거의 전제가 비로소 실제로 성립한다.

## 게이트 표 (변경 후)

| 필드 | 분류 | 출처 / 검증 |
|---|---|---|
| who · what | **필수(AI)** | 본문 앵커 (Phase009 L0/L3 게이트) |
| **published_at** | **기계 자동충전** | JSON-LD `datePublished`/OG — AI·앵커 불필요 |
| when(사건일) | 선택(AI) | 본문에 적혔으면 앵커, 없으면 생략. 유효앵커 0 → REVIEW |
| where · how · why | 선택(AI) | 앵커 또는 nil. 유효앵커 0 → REVIEW |

## 결정론적 성격

A·B 모두 입력(event6 + 원문/구조화 데이터)만으로 결정되는 순수 규칙이다. A는 필수/선택 분류
재배치라 비용 0. B는 이미 존재하는 추출기를 호출 경로에 잇는 것으로, 발행일 출처가 메타데이터라
판정 신뢰도 100%. AI 자기평가 없음.

## 열린 결정

- **extract.Apply 적용 지점**: 시드(`run`) 단계가 정석(처리 전 SKIPPED를 걸러 TODO 오염 방지)이나,
  최소 변경은 submit. 시드 적용 시 기존 세션 재시드 필요.
- **datePublished 부재 사이트**(OG만 있고 `published_time` 없음) 폴백: WARC 캡처일(근사, 출처 표기)
  vs `published_at` 빈칸 허용.
- **사건일(event date)**: 본문에 명시된 사건 발생일을 선택 `when`으로 받을지. 권장 — 본문에 적혔으면
  선택 앵커로, 없으면 생략(발행일은 `published_at`가 별도 보유).
- **기존 PASS 재판정**: when 앵커로 통과한 기존 레코드를 새 기준으로 다시 볼지(현 테스트 세션은 폐기됨).

## 다음 단계

- 구현: `gate.go` 필수/선택 재배치(A) + `cmd/submit`(·`run`/`next`)에 `extract.Apply` 배선(B)
  → impl → tsma 완수 → filefunc 0/0.
- 회귀 측정: when 제거 후 PASS율·REVIEW율 변화, `published_at` 채움률, Haiku/Sonnet 재실행으로
  when 헛앵커 0 확인.
- 후속(선택): Phase009에서 보류한 **when 앵커 타입 게이트**(사건일 한정 숫자/시간키워드 필수)는 본
  변경 후 필요성 재평가 — 발행일을 분리하면 사건일 헛앵커 위험이 크게 줄어 불필요할 수 있음.
