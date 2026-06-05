# Phase 005 — textmatch: 본문 포함 검증기 (규칙 헬퍼)

- 프로젝트: reins (GitHub: `park-jun-woo/reins`)
- 작성일: 2026-06-05
- 상태: 스캐폴드 빌드됨 (pkg/textmatch)
- 의존: Phase001(헌장)
- 패키지: `github.com/park-jun-woo/reins/pkg/textmatch` (외부 의존: `golang.org/x/text/unicode/norm`)

## 목적

"이 토큰이 원천에 **글자 그대로 실재하는가**"를 결정론으로 판정하는 순수 검증기. 게이트 규칙
패밀리 ④(치즈 방어)가 환각 차단에 *호출*하는 보편 헬퍼. ccnews `internal/anchor`의 핵심을 도메인
비결합으로 추출하고, 3-모델 테스트에서 드러난 **유니코드 false-negative를 NFC 정규화로 수정**한다.

핵심 명제 한 줄: *"정규화 후 substring. 같은 (source, token) → 같은 판정. 의미 추론 없음."*

## 배경 (수정 동기)

Opus/Sonnet/Haiku 30건 테스트에서 **세 모델 모두** 동일 false-negative 보고: 화면상 동일한 벵골어
결합모음(`বুড়ো`)·보스니아어 `đ`가 substring 실패 → 정직한 추출이 FAIL/DONE. 원인: ccnews
`normalize()`가 **공백만 접고 유니코드 정규화(NFC/NFD)를 안 해서** 합성형(é=`e`+`́` vs `é`)이 어긋남.
한 곳(검증기)에서 고치면 모든 소비 규칙이 혜택.

## API (순수 함수)

```go
package textmatch

// Normalize: 유니코드 NFC 적용 + 공백 런을 단일 스페이스로 접고 트림.
// 원천·토큰 양쪽 동일 적용(대칭 표면형 비교, 매핑 추론 아님).
func Normalize(s string) string

// Contains: token이 source의 substring인지(둘 다 Normalize 후).
// 빈/공백 토큰은 false (Contains(src,"")의 자명-참 함정 차단 — ccnews Phase009 L0 교훈 내장).
func Contains(source, token string) bool

// MissingTokens: tokens 중 source에 없는 것들(앵커 다건 검사 — 규칙이 Fact 만들기 좋게).
func MissingTokens(source string, tokens []string) []string
```

## 결정론 규칙

| 규칙 | 내용 |
|---|---|
| 정규화 | NFC → `strings.Fields` 공백 접기 → join. 원천·토큰 동일 적용. |
| 빈 토큰 | `Normalize(token)==""` → false. |
| 대칭성 | 한쪽만 변형 금지(매핑 추론 아님). |
| 의미 0 | 번역·동의어·퍼지 없음. 표면형 substring만. |

NFC 추가가 핵심 변경: 벵골어/분음부호 합성형 통일 → false-negative 소멸.

## 게이트 규칙에서의 사용 (Phase003 ④)

```go
//rule: id=who-anchor-present level=FAIL desc="필수 who 앵커가 원문에 실재해야"
var whoAnchorPresent = gate.Rule{
    Meta: gate.RuleMeta{ID: "who-anchor-present", Level: gate.LevelFail, Desc: "필수 who 앵커가 원문에 실재"},
    Check: func(ctx gate.Context) (bool, quest.Fact) {
        sub := ctx.Submission.(*Event6)
        if miss := textmatch.MissingTokens(ctx.Source, validAnchors(sub.Who.Anchors)); len(miss) > 0 {
            return true, quest.Fact{Where: "who.anchors", Expected: "원문 substring", Actual: miss[0]}
        }
        return false, quest.Fact{}
    },
}
```

규칙은 substring 로직을 재구현하지 않고 textmatch를 호출 — NFC 수정이 전 규칙에 일괄 적용.

## ccnews 추출 경로

- `ccnews/internal/anchor/normalize.go` → `textmatch.Normalize` + `norm.NFC` 선적용.
- `check_field.go`의 substring 루프 → `textmatch.Contains`/`MissingTokens`.
- ccnews `internal/anchor`는 verdict 매핑·필수/선택 분류·value 위생을 **게이트 규칙**으로 옮기고,
  substring 기계는 이 검증기에 위임(얇은 소비자).

## 결정론적 성격

(source, token)만으로 결정되는 순수 함수. IO·난수·외부상태 없음.

## 열린 결정

- **대소문자**: case-sensitive 유지(원어 표면형 충실) 권장.
- **분음부호 폴딩**: NFC(합성형 통일, 기본) vs NFKC(호환분해 — 과합치 위험으로 미채택).
- `MissingTokens` 외 다건 편의 API 표면.

## 다음 단계

- 구현 완료: `pkg/textmatch` + 단위테스트(`contains_test.go` — NFC 합성형: 벵골어 결합모음·`đ`·움라우트 합성/분해 동치).
- 회귀: ccnews가 textmatch로 전환 후, 3-모델 테스트의 벵골어/디아크리틱 FAIL이 PASS 회복되는지 재측정.
- Phase006(temporal)이 성분-앵커 숫자 검사에 textmatch 재사용.
