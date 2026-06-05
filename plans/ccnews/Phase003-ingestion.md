# Phase 003 — 투트랙 인제스천 (WARC 자동 다운로드)

- 프로젝트: ccnews · 작성일: 2026-06-03 · 상태: 완료
- 의존: Phase002(세션 스키마)

## 목적

CC-NEWS 전체를 두 방향으로 계속 빨아들이는 인제스천 엔진. WARC를 **CLI가 직접 다운로드**해
파싱하고 `articles`를 채운다. 어디까지 받았/처리했는지는 `ingestion` 커서로 추적.

## 투트랙

| 트랙 | 방향 | 동작 |
|---|---|---|
| **forward** | 최신 → | 최신 덤프부터 처리. 따라잡으면 `state=waiting`(새 덤프 폴링), 나오면 재개. |
| **backward** | → 과거 | 과거로 내려가며 모든 과거 뉴스 분석할 때까지 고반복. 최古 도달 시 `exhausted`. |

- WARC는 **최신부터 차례로** 받음. forward는 신규 끝, backward는 과거 끝에서 진행.
- 두 커서 독립. 동일 WARC 중복 처리 방지(`processed_warcs` 잠금).

**스케일 현실(사용자 확정 = "전체"):** CC-NEWS는 2016-08~현재로 WARC 수만~수십만 개·총 페타바이트급.
"전체 backward"는 **종료 목표가 아니라 무기한 진행 트랙**으로 설계한다 — 래칫이므로 언제 중단해도
커서에서 재개. 따라서 전량 완주를 전제로 한 가정(전체 디스크 동시 보관 등)은 두지 않는다:
WARC는 스트리밍 처리 후 보존 규칙(아래 열린결정)에 따라 회수. 데모는 부분 진행분으로 성립.

## WARC 파이프라인

```
Common Crawl CC-NEWS warc.paths (월별 목록)
   → 커서가 다음 WARC 선택 → 다운로드(.warc.gz)
   → 레코드 순회(response 레코드만): Target-URI=url, host, HTML payload, offset
   → 각 레코드 → articles[]에 TODO 추가 (warc{file,offset} 로케이터)
   → WARC 처리 끝 → processed_warcs에 추가(잠금), 커서 전진
```

## 결정론적 게이트
- **WARC 잠금**: `processed_warcs`에 있으면 재처리 안 함(ratchet — 중단 시 커서에서 재개).
- 다운로드 무결성: 크기/해시 확인(부분 다운로드 거부).

## CLI
- `ccnews run [--track forward|backward|both]` — 인제스천 루프 구동.
- 내부: 다운로드 → scan(레코드→기사) → (이후 Phase004~006 게이트가 기사 처리).

## 열린 결정
- **backward 시작점**: CC-NEWS 최古 ≈ 2016-08. 월별 `warc.paths` 획득 방식(CC 인덱스 HTTPS/S3).
- **forward 폴링 주기**: `waiting` 상태에서 새 덤프 확인 간격.
- WARC 로케이터: `{file, offset}` 직접 vs CDX 인덱스. (1차: file+offset)
- 다운로드 보관: 본문 미저장이라 앵커 재검증 입력이 WARC다. 따라서 **TODO 기사가 남은 WARC는
  보존**하고, 그 WARC의 기사가 전부 종단상태(PASS/REVIEW/DONE/BLOCKED/SKIPPED)가 된 뒤에만 삭제 가능.
  (삭제 후엔 그 기사 재처리 불가 — 종단상태라 NextTODO가 안 집으므로 정합.)

## 다음 단계
Phase004(robots 게이트) — 기사 처리 전 host robots 확인.
