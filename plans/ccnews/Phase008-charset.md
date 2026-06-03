# Phase 008 — 문자셋 감지 & UTF-8 정규화

- 프로젝트: ccnews (GitHub: `park-jun-woo/quest-ccnews`)
- 작성일: 2026-06-03
- 상태: 설계 — 구현 대기
- 의존: Phase003(인제스천/ReadBody), Phase005(추출), Phase006(앵커)

## 목적

비-UTF-8 인코딩(Shift-JIS·EUC-JP·EUC-KR·GB2312·windows-1251 등) 페이지를 **UTF-8로
디코드**한 뒤 추출·앵커를 수행해, 산출 value와 원어 anchors가 깨지지 않게 한다.

## 배경 (Phase007 찍먹 실측 근거)

실제 CC-NEWS WARC 1개 앞부분 10개 기사로 추출 캐스케이드를 검증한 결과:
- 캐스케이드 정상: **jsonld 9 / og 1 / SKIP 0**.
- 그러나 **일본 livedoor 기사 제목이 `��JR...` 모지바케**. 페이지가 Shift-JIS/EUC-JP인데
  UTF-8로 읽어 깨짐. CC-NEWS는 전 세계 뉴스라 레거시 인코딩이 상당수.

영향 범위: 깨진 본문에서는 JSON-LD 값/제목/event6 value가 쓰레기가 된다. **앵커는 같은
깨진 본문 기준이라 substring 검증 자체는 통과**하므로 환각 게이트는 안 깨지지만, **산출
품질**이 무너진다(원어 anchors·영어 value 모두 오염).

## 설계

### 디코드 지점 = `ingest.ReadBody` (단일 지점)
본문 바이트가 세션 밖 WARC에서 들어오는 유일한 통로가 `ReadBody`다. 여기서 **인코딩을
결정해 UTF-8로 변환한 바이트를 반환**한다. 다운스트림(`extract.Parse`, 앵커 대상 텍스트,
`next` 출력)은 항상 UTF-8을 받는다. → 한 곳만 고치면 전 경로가 UTF-8.

### 인코딩 결정 (결정론적 — 라이브러리에 위임)
`golang.org/x/net/html/charset.DetermineEncoding(content, contentType)` **한 호출**로 결정한다.
이 함수가 내부적으로 ① 선언된 Content-Type의 `charset=`, ② BOM, ③ HTML 앞 1024바이트
`<meta charset>`/`<meta http-equiv>` prescan, ④ 휴리스틱을 **이미 그 우선순위로 수행**한다
(x/net v0.27.0 확인). 따라서 우리가 Content-Type/meta/BOM을 각각 다시 구현하지 않는다 —
중복·과설계 회피.

- 입력: WARC HTTP 응답의 `Content-Type` 헤더(ReadBody가 추출) + body 바이트.
- `DetermineEncoding`이 `certain=false`로 자신 없게 판정하더라도 반환된 encoding으로 변환하되,
  변환 결과가 UTF-8로 무의미하면(아래 실패 정책) **UTF-8 가정**(현행 동작)으로 안전 폴백.
  (UTF-8 입력은 멱등 — DetermineEncoding이 utf-8을 돌려주므로 그대로 통과.)

### 구현
- 순수 변환 함수 분리(테스트 용이): `ToUTF8(raw []byte, contentType string) ([]byte, encName string)`.
  - `x/net/html/charset.DetermineEncoding`로 인코딩 판정 → `transform.NewReader`로 UTF-8 디코드.
  - 이미 UTF-8이면 그대로 반환(멱등).
- `ReadBody`는 HTTP response의 `Content-Type` 헤더를 추출해 `ToUTF8`에 넘긴다(현재 body 바이트만
  반환 → Content-Type을 함께 읽어 디코드 후 UTF-8 바이트 반환).
- 의존성: `golang.org/x/net`(직접 의존, v0.27.0) + `golang.org/x/text`(transform/encoding). `x/text`는
  이미 모듈 그래프에 v0.16.0으로 들어와 있어(x/net의 전이 의존) `import`만 추가하면 됨 — 새 외부 의존 도입 아님.

## 결정론적 성격

인코딩 결정·변환은 **결정론적 전처리**(같은 입력 → 같은 UTF-8 출력). 게이트가 아니라 게이트의
입력 품질을 보장하는 단계. 앵커 게이트(Phase006)는 이 UTF-8 텍스트에 적용되어야 하며, 앵커는
**디코드된 UTF-8 원어 표면형**이어야 한다(에이전트가 `next`에서 보는 본문과 동일 인코딩).

## 영향

- extract 품질↑: 일본/중국/러시아/한국 레거시 인코딩 기사의 title·body·event6가 올바른 원어로.
- 환각 게이트: 영향 없음(디코드는 next 표시·submit 검증 양쪽에 동일 적용 → 일관).
- (선택) `session.Extracted`에 `charset` 기록 — 디버깅/투명성용. 필수 아님.

## 열린 결정

- **디코드 실패/미지원 charset 정책**: UTF-8 가정으로 진행(현행) vs 품질 플래그/`SKIP`. 기본은
  진행(과도한 SKIP 회피). 단 변환 후에도 invalid UTF-8 비율이 높으면 로깅 정도.
- `ReadBody` 시그니처 변경 방식: 내부에서 디코드해 UTF-8 바이트만 반환(다운스트림 무변경, 권장)
  vs `(bytes, contentType)` 반환 후 호출측 디코드.
- `Extracted.charset` 기록 여부.

## 다음 단계

- 구현: `ToUTF8` 추가(순수) + `ReadBody` 디코드 연동 → impl → tsma 완수 → filefunc 0/0.
- 재검증: Phase007 찍먹 재실행해 livedoor 제목이 정상 일본어로 나오는지 확인.
