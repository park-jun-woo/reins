# Phase 006 — temporal: 시간 명세 정규화 검증기 (규칙 헬퍼)

- 프로젝트: reins (GitHub: `park-jun-woo/reins`)
- 작성일: 2026-06-05
- 상태: 스캐폴드 빌드됨 (pkg/temporal — v1: gregorian+interval+relative)
- 의존: Phase001(헌장), Phase005(textmatch — 성분-앵커 검사 재사용)
- 패키지: `github.com/park-jun-woo/reins/pkg/temporal`

## 목적

AI가 산출한 **구조화 시간 명세**(역법/성분/오프셋)를 받아 **정규 그레고리력 ISO**(단일·기간)로
결정론 변환·검증하는 순수 검증기. 게이트 규칙 패밀리 ①(형식)·④(사실)가 *호출*한다. 다양한 입력
형태(절대-다역법 / 기간 / 상대표현)를 허용하되, 다국어 *읽기*는 AI에 두고 검증기는 **변환·산술만**.

핵심 명제 한 줄: *"AI는 어느 역법·성분·오프셋인지 식별, 기계는 그걸 그레고리력 ISO로 변환·검증.
역법 변환·날짜 산술은 언어무관 결정론."*

## 배경

3-모델 테스트에서 `when`이 가장 약했다: 날짜가 본문 아닌 메타에 살거나(→ ccnews Phase010
published_at로 해소), 상대표현·비그레고리력이라 단일 ISO로 환원 불가. `Spring 2026`·`2026-06-11 to
2026-07-19`·`1395-10-21 (Persian)`·`This Wednesday` 등 — **역법·기간·상대는 전부 결정론 변환 가능**.

## 입력 스키마 (규칙이 채울 spec; value는 검증기가 산출)

```jsonc
{
  "kind": "absolute" | "relative",
  "calendar": "gregorian|persian|islamic|hebrew|chinese",  // (absolute) AI 식별
  "start": "1395-10-21",         // 해당 역법 성분
  "end":   "1395-10-25" | null,  // null=시점, 채우면 기간
  "offset_days": -1,             // (relative) 오늘=0·어제=-1·그제=-2·내일=+1
  "anchors": ["۲۱ دی ۱۳۹۵"]      // 원문 토큰 — textmatch가 본문 substring 검증(규칙 ④)
}
```

## 세 입력 모드 (전부 기계가 정규화)

| 모드 | AI가 주는 것 | 검증기 결정론 처리 | 출력 value |
|---|---|---|---|
| 절대(역법) | calendar + start | 역법→그레고리력 변환 | `2017-01-10` |
| 기간 | start + end | 양끝 변환 | `2017-01-10/2017-01-14` |
| 상대 | offset_days | ref(published_at) + offset | `2026-06-02` |

## API (순수 함수)

```go
package temporal

type Kind string                              // "absolute" | "relative"
const ( Absolute Kind = "absolute"; Relative Kind = "relative" )

type Calendar string                          // v1은 Gregorian만 변환
const ( Gregorian Calendar = "gregorian"; Persian; Islamic; Hebrew; Chinese )

// Spec: AI 읽기로 채운 구조화 시간 명세(value는 검증기가 산출).
type Spec struct {
    Kind       Kind
    Calendar   Calendar
    Start      string   // 해당 역법 성분
    End        string   // "" = 시점, 채우면 기간
    OffsetDays int      // (relative) 오늘=0·어제=-1·내일=+1
    Anchors    []string // 원문 토큰 — ComponentsInAnchor가 검증
}

type Result struct { Value string; IsInterval bool; Determined bool }

// Resolve: spec → 그레고리력 ISO. 상대는 ref(=ref.AddDate). 비그레고리 역법·파싱 실패 →
// Determined=false (규칙이 REVIEW 레벨로 매핑). v1은 표준 time만 사용.
func Resolve(spec Spec, ref time.Time) Result

// ComponentsInAnchor: start/end의 라틴숫자 런 성분이 anchor에 실재하는지(textmatch.Contains 재사용).
// 날짜를 엉뚱한 토큰에 건 케이스(Türkiye·Kurtulmuş)를 when에 한해 차단. v1은 라틴숫자만(타 숫자체계 v2).
func ComponentsInAnchor(spec Spec, anchors []string) bool
```

## 게이트 규칙에서의 사용 (Phase003 ①·④)

```go
//rule: id=when-format level=FAIL desc="when 시간 명세가 정규형으로 변환 가능해야"
//rule: id=when-anchor-support level=FAIL desc="when 성분이 원문 앵커에 실재해야(엉뚱앵커 차단)"
//rule: id=when-undetermined level=REVIEW desc="모호·기준일부재 날짜는 사람 확인"
```

`Resolve(...).Determined==false` → REVIEW 규칙 발동. `!ComponentsInAnchor(...)` → FAIL 규칙 발동.

## 결정론 경계 (정직하게)

- **결정론**: 역법 변환·날짜 산술·숫자체계 매핑(۱۲۳→123 유한표) — 같은 입력 같은 출력, 언어무관.
  (v1 실구현은 라틴숫자 성분 매칭·그레고리 변환·상대 산술까지; 역법 변환·비라틴 숫자 매핑은 v2.)
- **AI-신뢰 잔여**: `calendar`·`offset_days`는 AI의 *읽기*. 다국어 표 없이 검증 안 함(Phase009에서
  거부한 "다국어 스톱워드 군비경쟁" 함정). 대신 anchors 실재 검증 + offset은 작은 정수(−2~+1)라 위험
  작음. 잔여로 명시.
- **Undetermined**: 미지원 역법·상대인데 ref 부재·모호(`Spring 2026`·"soon") → REVIEW. 본질적 미정.

## 범위 슬라이싱 (한 번에 다 안 함)

| 단계 | 범위 | 근거 |
|---|---|---|
| **v1** | gregorian + interval + 상대(ref 기준) | 뉴스 대다수 커버, 외부 의존 최소(표준 time만) |
| v2(조건부) | persian/islamic/hebrew/chinese 변환 라이브러리 | 해당 언어 기사 비중 의미 있을 때만 — "관찰 안 된 건 선제 구현 안 함"(Phase009 원칙) |

v1은 비그레고리 역법을 `Determined=false`(REVIEW)로 정직 반환, v2에서 점증.

## 결정론적 성격

`Resolve`/`ComponentsInAnchor`는 (spec, ref) 또는 (spec, anchors)만으로 결정되는 순수 함수.

## 열린 결정

- v1 지원 숫자체계(라틴 + 페르시아·아랍·데바나가리·한자 — 유한표).
- interval 표기(ISO `start/end`).
- 상대 ref 출처(소비자 주입; ccnews=published_at).
- offset 고빈도 교차검사(선택; 기본 미채택 — AI 신뢰).
- v2 역법 라이브러리 선정(검증·라이선스).

## 다음 단계

- 구현 완료: `pkg/temporal` v1(gregorian+interval+relative) + 단위테스트(`temporal_test.go` — 상대 산술·interval·undetermined·성분앵커).
- 통합(ccnews): `when`을 본 spec로 받아 `Resolve(spec, published_at)`로 정규화, 게이트 규칙으로 검증.
- 재측정: when 헛앵커(Türkiye·Kurtulmuş류)가 `ComponentsInAnchor` 규칙으로 차단되는지 3-모델 재실행.
