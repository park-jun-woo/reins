# Phase 002 — 세션 스키마 & 상태기계

- 프로젝트: ccnews · 작성일: 2026-06-03 · 상태: 설계 확정
- 의존: Phase001(개요)

## 목적

런타임 상태 SSOT인 `session.json` 스키마와 기사 퀘스트 상태기계, `internal/session` Go
구조체를 정의한다. (생성=AI, 판정=기계: 세션은 *판정 결과*만 담고 본문은 안 담는다.)

## 세션 JSON 스키마

```jsonc
{
  "version": 1,
  "user_agent": "parkjunwoo-quest/0.1 (+https://www.parkjunwoo.com)",

  // 투트랙 인제스천 커서 (상세 Phase003)
  "ingestion": {
    "source": "cc-news",
    "forward":  { "cursor": "2026-06/CC-NEWS-...warc.gz", "state": "running" },   // running|waiting
    "backward": { "cursor": "2026-05/CC-NEWS-...warc.gz", "state": "running" },   // running|exhausted
    "processed_warcs": ["..."]
  },

  // 호스트 캐시 — robots + 매체정보. (추출 템플릿 없음: 구조화 데이터만 신뢰)
  "hosts": {
    "www.example.com": {
      "media_name": "Example News",          // JSON-LD publisher.name / og:site_name
      "site_url": "https://www.example.com", // 매체사이트주소 (사용자 명시 수집항목)
      "robots": { "fetched_at":"...", "robots_url":"https://www.example.com/robots.txt",
                  "status":"ok", "crawl_allowed":true, "crawl_delay_sec":1 },   // status: ok|missing(=허용)|unreachable; crawl_allowed = 크롤링허용여부(사용자 명시)
      "license": null                        // CC저작권종류: {type,url,source} | null
    }
  },

  // 기사 퀘스트 (work-list)
  "articles": [
    {
      "warc": { "file":"CC-NEWS-...warc.gz", "offset":123456 },
      "url": "https://www.example.com/news/123",   // WARC Target-URI
      "host": "www.example.com",
      "lang": "ko",                          // 원문 언어코드 (산출은 영어)
      "state": "TODO",                       // 아래 상태기계
      "skip_reason": null,                   // BLOCKED/SKIPPED 시 사유
      "tries": 0,
      "extracted": { "title":"...", "author":"...", "published_at":"2026-06-02",
                     "source":"jsonld", "body_len":3412 },   // source: jsonld|og. 본문 텍스트 미저장
      "event6": {                            // value=영어, anchors=원어 (Phase006)
        "who":  { "value":"Lee Jae-yong, Samsung Chairman", "anchors":["이재용","삼성전자"], "anchored":true },
        "when": { "value":"2026-06-02", "anchors":["6월 2일"], "anchored":true },
        "what": { "value":"...", "anchors":["..."], "anchored":true },
        "where": null, "how": null, "why": null
      },
      "verdict":"PASS", "verdict_reason":"...", "collected_at":"...", "log":[]
    }
  ]
}
```

## 상태기계 (단방향)

```
TODO ─┬─► PASS      필수(who/when/what) 앵커 완료
      ├─► REVIEW    해석필드(how/why) 사람 확인 필요
      ├─► DONE      MaxTries 초과(추출/앵커 실패)
      ├─► BLOCKED   robots 거부 (skip_reason 기록)
      └─► SKIPPED   구조화 데이터 없음 — 신뢰 불가 (skip_reason 기록)
```
PASS/REVIEW/DONE/BLOCKED/SKIPPED는 잠금(불가역). NextTODO는 TODO만 집는다.

## Go 구조체 (internal/session, 잠정)
- `Session{ Version, UserAgent, Ingestion, Hosts map[string]*Host, Articles []*Article }`
- `Host{ MediaName, Robots, License }`
- `Article{ WARC, URL, Host, Lang, State, SkipReason, Tries, Extracted, Event6, Verdict, ... }`
- `Event6{ Who,When,Where,What,How,Why *Field }`, `Field{ Value, Anchors []string, Anchored }`
  (`Interpretive`는 표시 라벨일 뿐 게이트 입력 아님 — REVIEW는 `len(Anchors)==0`로 기계 판정. Phase006)
- load/save(JSON), NextTODO, Find, Counts (comail/session 패턴 재사용)

## 설계 결정
1. **본문 미저장(파생 결정)** — 사용자가 명시한 "사실만 농축·공개 안전"에서 따라온 표현 복제
   회피책. `session.json`에는 `extracted.body_len`만 남기고 본문 텍스트는 직렬화하지 않는다.
   본문은 **WARC가 원천**이므로, 재처리·재검증이 필요하면 `warc{file,offset}`로 다시 읽는다
   (앵커 검증 입력은 세션이 아니라 WARC). 따라서 "미저장"이 검증 재현성을 깨지 않는다.
2. **템플릿 레지스트리 없음** — 구조화 데이터만 신뢰하므로 호스트 캐시는 robots+매체정보만.
3. **skip_reason 단일 필드** — BLOCKED/SKIPPED 공용, JSONL audit 레코드로 emit(Phase007).

## 열린 결정
- `Field.anchors` 빈 배열 허용 정책: present한 선택필드가 `anchors=[]`면 REVIEW로 자동 분류
  (필수필드가 `anchors=[]`면 FAIL). 구조적 트리거 — 주관 판정 없음. Phase006 verdict 표 참조.

## 다음 단계
Phase003(인제스천) — `ingestion` 커서로 WARC 받아 `articles` 채우는 scan/run.
