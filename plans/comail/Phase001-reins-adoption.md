# Phase 001 — comail의 reins 이식 (첫 소비자 찍먹)

- 프로젝트: comail (GitHub: `park-jun-woo/quest-comail`)
- 작성일: 2026-06-05
- 상태: 구현 완료 (Phase002 범위 빌드·테스트·tsma·filefunc green) · G5/export는 Phase007로 이월

## 목적

comail(회사명 → 이메일 수집 퀘스트)을 **reins 위로 이식**한다. 자체
`internal/session`·`internal/gate`·`cmd/*`를 버리고, 도메인 로직을 `gate.Definition`
**4메서드**(Seed/Render/Prepare/Rules)로만 재구현한다. 래칫·명령 골격·레벨 집계·export는
reins가 공급한다.

핵심 명제: *"comail은 reins의 첫 실사용 소비자다. 과제가 단순해(이메일 1개 = 1퀘스트)
프레임워크의 추상화가 진짜 도메인을 견디는지 가장 싸게 검증한다 — N=1 추상화 금지의 N=1."*

이 Phase의 산출물은 **동작하는 이식**이 아니라 **이식 설계 + reins에 뚫어야 할 구멍 목록**이다.
comail을 끼워보며 드러난 프레임워크 갭이 곧 reins v1 안정화 입력이 된다.

## 현 comail ↔ reins 매핑

| 현 comail | reins | 처리 |
|---|---|---|
| `internal/session.Quest{Company,State,Email,Site,Tries,Log}` | `quest.Item{Key,State,Tries,Payload,Log,…}`(+Verdict/Reason/CollectedAt/Emitted) | **삭제**. Company→`Key`, Email/Site→`Payload`로 |
| `internal/session.State` (TODO/PASS/REVIEW/DONE) | `quest.State` (+SKIPPED/BLOCKED) | **삭제**. reins 것 사용(여분 상태는 무시) |
| `internal/session.MaxTries = 3` | `quest.MaxTries = 3` | 값 일치 — 그대로 |
| `session.Save/Load/New/NextTODO/Find` | `quest.Session` + `cli` 로더 | **삭제**. reins 공급 |
| `gate.Verify(Candidate) (Verdict, reason)` — 순차 | `[]gate.Rule` + `gate.Evaluate` — 카탈로그 | **재작성**(아래 §게이트) |
| `gate.Candidate{Company,Site,Email,Source}` | `Submission`(Prepare가 디코드) | 유지하되 **`mxResolved bool` 1필드 추가**(Prepare가 박제, §부작용 격리) |
| `cmd/{scan,next,submit,status,export}` | `cli.NewQuestCmd` | **삭제**. 4메서드만 남김 |
| `cmd/export_xlsx,export_csv` | (reins엔 없음) | **갭 — §열린 결정 4** |

이식 후 comail은 `main.go` 한 줄 + `Definition` 구현 1개 + 규칙 카탈로그로 줄어든다:

```go
func main() { cli.NewQuestCmd("comail", comailDef{}, cli.Options{Version: "0.2"}).Execute() }
```

## 게이트: 순차 `Verify` → 규칙 카탈로그 (이식의 핵심)

현 `Verify`는 **이른 반환(early-return)** 으로 우선순위를 표현한다. reins는 **모든 규칙을 돌리고
레벨로 집계**한다(any Fail→FAIL, else any Review→REVIEW, else PASS). 순서 의존을 규칙 독립성으로
번역해야 한다. (reins `Level`은 **Fail/Review 2종뿐** — Pass-level 규칙은 없다. "위반 없음=PASS"라
분해표에 Pass 규칙이 없는 게 정상이다.)

### 분해표 (현 순차 검증 → 규칙 8개 + Prepare short 1개)

| # | 현 Verify 분기 | reins 규칙 ID | Level | Check 발동 조건 |
|---|---|---|---|---|
| 1 | 이메일 형식 오류 | `email-format` | Fail | `!emailRe.Match(email)` |
| 2 | 블랙리스트(플레이스홀더) | `placeholder-domain` | Fail | `blacklist[edom]` |
| 3 | 프리메일 | `freemail-domain` | Review | `freemail[edom]` |
| 4a | 사이트 URL 누락 | `site-missing` | Fail | `regDomain(site) == ""` |
| 4b | 불일치 + 출처도 비공식 | `mismatch-unofficial-source` | Fail | `edom≠sdom && srcdom≠sdom` |
| 5a | 출처 fetch 실패 | (규칙 아님 → **Prepare short**) | Fail | fetch 에러 — §부작용 격리 |
| 5b | 출처에 이메일 없음(환각) | `source-lacks-email` | Fail | `MissingTokens(ctx.Source,[email])≠∅` |
| 6 | 불일치(그룹사 공용 추정) | `group-domain-review` | Review | `edom≠sdom && srcdom==sdom` |
| 7 | MX 없음 | `mx-missing` | Fail | `!mxResolved`(Prepare가 조회) |

> 현 `sourceHasEmail`(source_has_email.go)은 **fetch 실패와 문자열 부재를 한 FAIL로 묶는다**. 이식은
> 이를 둘로 쪼갠다 — fetch(부작용)는 `Prepare`의 short verdict(5a), 문자열 대조(순수)는 규칙 5b.
> 그래서 "현 7단계 → 규칙 8개 + short 1개"가 정확한 카운트다.

### 레벨 집계가 순서를 보존하는가 — 검증

집계는 **Fail > Review > Pass** 이므로, "Fail 분기가 Review 분기보다 먼저 반환"하던 순서는 자동 보존된다:

- 환각(5b, Fail)이 그룹사-REVIEW(6, Review)보다 먼저 → 둘 다 발동해도 **Fail 우선** = 순차와 동일 ✅
- 블랙리스트(2, Fail)가 프리메일(3, Review)보다 우선 → 동일 ✅

순서가 **보존 안 되는** 케이스(=의미 분기점, §열린 결정 1·2):

- **프리메일에 source-lacks-email을 적용하나?** 현 코드는 프리메일이면 step3에서 **즉시 REVIEW
  반환** → 출처 재확인 안 함. 카탈로그에선 `freemail-domain`(Review)과 `source-lacks-email`(Fail)이
  동시 평가 → 프리메일인데 출처에 그 메일이 없으면 **FAIL로 뒤집힘**. (논쟁의 여지: 더 엄격해짐)
- **그룹사 공용 도메인에 MX를 검사하나?** 현 코드는 step6에서 REVIEW 반환 → step7(MX) 도달 안 함.
  카탈로그에선 `mx-missing`이 edom에 대해 항상 평가 → 그룹 도메인 MX 불량이면 REVIEW→FAIL.

→ **방어적 규칙 작성 원칙**: 각 규칙은 입력을 독립 파싱하고, 평가 불가한 입력엔 *no-fire*
한다(예: `email-format`이 이미 잡을 잘못된 이메일에 대해 다른 규칙은 발동 안 함). 그래야 Facts가
중복 노이즈로 차지 않는다.

## 부작용 격리: `Prepare`가 fetch·MX를 흡수

reins `gate.Context`의 `Source`는 **캐시된 본문(ground truth)** 이라는 게 도그마다(`Source string //
re-confirmed by cheese-defense rules`). 현 comail은 `gate.Verify` **안에서** 출처 URL을 HTTP fetch
하고 `net.LookupMX`를 호출한다 — 비결정·테스트 불가.

이식 원칙: **부작용은 전부 `Prepare`로 올린다. 규칙은 캐시된 Context에 대해 순수해진다.**

전제: `Candidate`에 `mxResolved bool`을 추가하고(§매핑표), `comailDef`와 규칙·`regDomain`/`fetch`는
**같은 패키지**에 둔다(현 `internal/gate`의 비공개 자산 `regDomain`/`fetch`/`emailRe`를 재호출하려면
패키지 경계를 넘지 않아야 함 — §열린 결정 5의 배치 결정이 이 코드의 선결조건).

```go
func (d comailDef) Prepare(it *quest.Item, raw []byte) (gate.Context, *quest.Verdict, error) {
    var sub Candidate // {Company,Site,Email,Source, mxResolved bool}
    if err := json.Unmarshal(raw, &sub); err != nil { return gate.Context{}, nil, err }

    body, err := fetch(sub.Source)               // ← 부작용 1: HTTP
    if err != nil {                              // 접근 실패 = 환각 검증 불능 → 즉시 FAIL(short)
        v := &quest.Verdict{Outcome: quest.OutFail,
            Facts: []quest.Fact{{Rule: "source-fetch", Where: "source",
                Expected: "fetchable page", Actual: err.Error()}}}
        return gate.Context{}, v, nil
    }
    if emailRe.MatchString(sub.Email) {             // 형식 통과 시에만 — no-fire 원칙(§아래)을 Prepare도 지킨다
        sub.mxResolved = lookupMX(regDomain(sub.Email)) // ← 부작용 2: DNS, 결과를 페이로드에 박제
    } // 형식 오류면 mxResolved=false로 두고 DNS 생략 — email-format 규칙이 FAIL을 내고 mx-missing은 가려진다

    return gate.Context{Item: it, Submission: &sub, Source: body}, nil, nil
}
```

- `source-lacks-email` 규칙 = `textmatch.MissingTokens(ctx.Source, []string{email})` — **순수**, reins
  `pkg/textmatch` 재사용. ⚠️ **회귀 주의**: 현 `sourceHasEmail`은 `strings.Contains(ToLower(body),
  email)`로 **대소문자 무시**지만, `textmatch.Normalize`는 NFC+공백접기만 하고 **case-fold 안 한다**
  (`pkg/textmatch`에 ToLower/Fold 0건). 그대로 갈아타면 `Contact@ACME.com` vs 페이지의
  `contact@acme.com`이 FAIL이 되는 **거짓양성 회귀**. → email/source를 ToLower 후 넘기거나 textmatch에
  case-fold 옵션을 신설(§열린 결정 6).
- `mx-missing` 규칙 = `!sub.mxResolved` 읽기 — **순수**.
- 결과: `gate.Evaluate(rules, ctx)`가 **(rules, ctx)만으로 결정적** — reins 도그마 충족, 규칙 단위
  테스트가 네트워크 없이 가능.

> **⚠ G5 — fetch 선행 부작용 (구현 후 스모크로 발견).** 이 격리는 부작용을 Prepare로 *모았을* 뿐
> *스케줄링*은 못 한다. 위 코드는 `fetch(sub.Source)`를 **이메일 형식검사보다 먼저** 실행한다
> (구현 `prepare.go`도 동일). 그래서 형식 틀린 입력도 네트워크를 태우고, fetch 실패 시 short
> verdict가 `source-fetch`로 나가 피드백 사유가 `email-format`이 아니라 엉뚱한 source-fetch로
> 뒤바뀐다(판정 FAIL은 동일, 근본원인만 오인). MX는 `emailRe` 가드 뒤로 뺐으나 **fetch엔 같은 가드가
> 없다.** 근본 해결은 규칙 우선순위(defeat) 기반 단계 평가 — `plans/reins/Phase007-toulmin-gate.md`
> 에서 처리하며, 그때 comail이 재이식(Phase012)된다.

## reins에 뚫어야 할 구멍 (프레임워크 갭)

comail을 끼우다 막힌 지점 — reins v1 안정화의 실입력:

| 갭 | 증상 | 후보 해법 |
|---|---|---|
| **G1 submit 입력** | reins `submit`은 `--key`+`--in <file>\|-`(raw 바이트). comail은 `--company/--site/--email/--source` 4플래그 | Render가 **JSON 제출 스키마**를 프롬프트로 명시 → 에이전트가 JSON을 stdin/파일로(readSubmission이 raw를 Prepare에 직통). **선결(Phase002): Render가 출력할 스키마의 필드명이 Prepare `json.Unmarshal`의 `Candidate` 필드(Site/Email/Source)와 1:1 일치해야 함.** (또는 reins에 도메인 플래그 주입 훅) |
| **G2 scan 인자** | reins `scan`은 위치인자만 Seed로 전달. comail `--col/--sheet` 자리 없음 | `comail scan companies.xlsx A Sheet1` 위치인자로 인코딩, Seed가 파싱 |
| **G3 원본 메타 보존** | 현 comail `session.Session`은 **`SourcePath/SourceSheet/SourceCol` 3필드를 1급으로 보유**(export가 원본 형식 복원에 사용). 그러나 reins `quest.Session{Version,Items}`엔 둘 곳이 없어 이식 시 이 메타가 갈 곳을 잃음 | **1순위: reins `quest.Session`에 `Meta map[string]any` 확장**(도메인이 갖던 1급 필드의 정직한 이전). 2순위(비추): 센티넬 `Item`(`Key:"__source__"`) — `Export`/`NextTODO`/`status` 집계에 가짜 아이템을 흘리는 **오염 우회**라 지양 |
| **G4 xlsx/csv export** | reins JSONL `Sink` 고정이 **두 곳**: `export` 명령(`newExportCmd`→`newJSONLSink`)뿐 아니라 **`submit`도 매 제출 시 강제 JSONL export**(submit.go:59 `newJSONLSink`). comail 간판기능인 원본 xlsx 복사+이메일/상태 열은 자리 없음 | (a) reins root에 **별도 `comail export` cobra 명령** 부착(단 submit의 강제 JSONL은 여전히 남음). (b) `cli.Options.Sink` 훅 신설 — **submit·export 양쪽**을 덮어야 함. (`Options.ExtraCommands`로 명령 확장도) |

G1·G2는 도메인 측 우회로 흡수 가능. **G3·G4는 reins 코어 확장 후보** — 두 번째 소비자(ccnews)도
원본 보존/커스텀 출력이 필요하면 `Options` 훅으로 승격한다(N=2에서 안정화 원칙).

## 열린 결정

1. **프리메일 × 출처 재확인 — [확정: 현 동작 보존]** 구현이 `rule_source_lacks_email.go`의
   `passedEarly` 가드로 프리메일을 source-lacks-email에서 배제 → 프리메일=무조건 REVIEW 유지.
   (§레벨 집계 검증의 "보존 안 되는 케이스"는 이 가드로 막힘.)
2. **그룹사 도메인 × MX — [확정: 현 동작 보존]** 구현이 `rule_mx_missing.go`의 `if p.mismatch`
   no-fire로 그룹 도메인 MX 검사 생략 → 현 동작 유지.
3. **submit UX(G1) — [확정: JSON stdin 제출]** `render.go`/`agent_prompt.go`가 `{site,email,source}`
   스키마로 `comail submit --key … --in -`를 출력. reins 무수정 경로.
4. **export(G4)** — comail 전용 명령으로 둘까, `Options` 훅을 reins에 뚫을까? ccnews는 JSONL이라
   당장은 comail만의 요구 → 일단 전용 명령, ccnews 합류 시 승격 검토.
5. **패키지 배치 (선결)** — `comailDef`·규칙·`emailRe`/`freemail`/`blacklist`/`regDomain`/`fetch`를
   어디 둘까? 비공개 자산을 규칙에서 호출하려면 **한 패키지**여야 한다(§부작용 격리 코드의 전제).
   `internal/`에 다 남길지, 범용 `regDomain`(eTLD+1)만 reins `pkg`로 승격할지.
6. **case-fold (선결)** — 이메일 대조 회귀(§부작용 격리). comail 측 ToLower 전처리로 막을지,
   reins `pkg/textmatch`에 case-fold 옵션을 신설할지. ccnews 앵커도 원어 대소문자 영향을 받으므로
   textmatch 확장이면 두 소비자 공통 이득.

## 다음 단계

- **Phase002 — Definition 구현 [완료]**: `comailDef`(Seed/Render/Prepare/Rules) + 규칙 8개, `go.mod`에
  `require github.com/park-jun-woo/reins`(+`replace => ../`) 추가 완료. 단위테스트는 규칙별 순수 테이블 테스트.
- **Phase003 — export 결정 반영**: G3/G4 해법 확정 후 xlsx write-back 재배선.
- 그 뒤 **실루프 검증**: `scan → next → (에이전트) → submit → status → export` 1회 종단 — reins의
  첫 end-to-end 실사용 통과를 comail로 찍는다.
- **Phase007 분기 (게이트 백엔드 교체)**: `plans/reins/Phase007-toulmin-gate.md`가 reins 게이트를
  toulmin defeat 그래프 + ground provider로 끌어올린다. 본 Phase의 8규칙+가드는 거기 **Phase012에서
  엣지 5~6개 그래프로 재이식**되며 증발 예정 — comail이 그 N=1 증명대다. (G5도 그때 근본 해소.)
