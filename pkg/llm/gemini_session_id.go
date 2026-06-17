//ff:func feature=llm type=helper control=sequence level=error
//ff:what geminiSessionID — geminiRandRead seam(표준 crypto/rand.Read)으로 RFC 4122 v4 UUID 문자열을 생성한다(외부 의존성 없음). gemini Continue 1차 호출에서 reins가 세션을 만들 때 --session-id로 넘길 식별자. 16바이트를 읽어 version(0x40)·variant(0x80) 비트를 박고 8-4-4-4-12 하이픈 포맷으로 찍는다. 리더 실패 시 에러(실 crypto/rand는 사실상 발생 안 함; 테스트는 seam으로 주입). codex/grok처럼 백엔드 내부 헬퍼.

package llm

import "fmt"

// geminiSessionID generates an RFC 4122 v4 UUID string using only crypto/rand (no
// external dependency). reins issues it as --session-id when it creates the session
// on the first Continue call. Version (0x40) and variant (0x80) bits are set per the
// v4 spec; a rand.Read failure (practically never) is returned as an error.
func geminiSessionID() (string, error) {
	var b [16]byte
	if _, err := geminiRandRead(b[:]); err != nil {
		return "", fmt.Errorf("gemini session id: %w", err)
	}
	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // variant 10
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16]), nil
}
