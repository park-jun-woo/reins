# Phase 005 — 구조화 데이터 추출 + 신뢰 게이트

- 프로젝트: ccnews · 작성일: 2026-06-03 · 상태: 완료
- 의존: Phase002(세션), Phase004(robots)

## 목적

WARC HTML에서 제목·본문·작성자·발행일을 **구조화 데이터만으로** 결정론적 추출한다.
**1단계로 끝낸다** — 셀렉터 학습도, AI DOM 추측도 없음. 구조화 데이터가 없는 사이트는
*self-declare하지 않은 = 신뢰할 수 없는* 것으로 보고 **SKIP**한다.

> 설계 변경 근거: JSON-LD/OG조차 없는 사이트의 본문을 AI/셀렉터로 추측해 뽑으면 신뢰가 안 됨.
> "구조화 데이터를 스스로 선언한 사이트만 신뢰"가 더 결정론적이고 방어적. AI는 추출에서 손 뗌.

## 추출 (결정론적 파스)

```
1) JSON-LD  <script type="application/ld+json"> @type=NewsArticle/Article
     headline → title, articleBody → body, author → author,
     datePublished → published_at, publisher.name → media_name, inLanguage → lang
2) OG/메타 폴백  og:title, og:site_name, article:published_time, article:author
3) 둘 다 없음/불충분 → SKIPPED
```

## 신뢰 게이트 (히트 판정)
- **PASS 조건**(둘 다): ① 구조화 출처(`source=jsonld|og`)에서 `title` 비어있지 않음
  (= 사이트가 self-declare함) + ② 앵커 대상 텍스트 `body_len ≥ N` 확보.
  `body_len`은 `articleBody`가 있으면 그 길이, 없으면 **HTML 태그제거 전문**의 길이(위 "앵커 대상 본문").
- **SKIPPED**: JSON-LD/OG의 Article 자체가 없음(self-declare 안 함) → `skip_reason="구조화 데이터 없음(JSON-LD/OG) — 신뢰 불가"`.
- **SKIPPED**: 구조화는 있으나 앵커 대상 텍스트가 `N`자 미만(앵커 검증 불가) → `skip_reason="본문 텍스트 부족 — 앵커 불가"`.
- 의견 없음 — 구조화 필드 존재 / 텍스트 길이만 본다.

## AI 역할
- 추출 단계엔 **AI 없음**(순수 파싱). AI는 다음 단계(Phase006 event6 생성)에서만.

## 앵커 대상 본문
- event6 앵커(Phase006)는 **원문 텍스트에 substring 존재**하는지만 본다.
- 앵커 대상 텍스트 우선순위: ① JSON-LD `articleBody`가 있으면 그것. ② 없으면 WARC HTML 태그제거 전문.
- **모순 아님 해명**: 사용자의 "구조화만 신뢰"는 *추출 필드*(title/author/published/media)를 무엇에서
  뽑느냐의 규칙이다. 앵커 게이트는 *추출이 아니라* "AI가 말한 토큰이 원문에 실재하나"의 대조이고,
  그 원문은 WARC HTML 전문이 정당한 원천이다(셀렉터로 의미 필드를 추측하는 게 아님).
  따라서 `articleBody` 부재 시에도 SKIP이 아니라, **구조화 추출 PASS면(아래) HTML 전문을 앵커 대상으로** 쓴다.
- 본문 텍스트 자체는 세션에 저장 안 함(WARC가 원천 — Phase002 결정 1).

## 열린 결정
- `body_len` 최소 N자 기준값(앵커 검증이 의미 있으려면 몇 자 이상인지).
- JSON-LD 다중/중첩(@graph) 파싱 우선순위.
- (해결됨) `articleBody` 필수 아님 — 부재 시 HTML 태그제거 전문을 앵커 대상으로. 위 신뢰 게이트 참조.

## 다음 단계
Phase006(event6 앵커 게이트) — 추출 본문에서 육하원칙 생성 + 원어 앵커 검증.
