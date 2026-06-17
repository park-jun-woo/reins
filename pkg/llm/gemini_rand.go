//ff:type feature=llm type=adapter
//ff:what geminiRandRead — geminiSessionID가 쓰는 난수 seam(패키지 변수). 본문은 표준 crypto/rand.Read를 가리킨다(외부 의존성 없음). execGemini와 같은 패턴의 per-함수 변수라, gemini 테스트가 실패하는 리더를 주입해 UUID 생성 에러 경로(사실상 발생 안 하는 rand.Read 실패)를 진짜 엔트로피 고갈 없이 검증할 수 있다.

package llm

import "crypto/rand"

// geminiRandRead is the randomness seam used by geminiSessionID. The body is the
// standard crypto/rand.Read (no external dependency); keeping it a package var (like
// execGemini) lets a test inject a failing reader to exercise the UUID-generation
// error path without exhausting real entropy.
var geminiRandRead = rand.Read
