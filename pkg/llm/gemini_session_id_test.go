//ff:func feature=llm type=helper control=sequence level=error
//ff:what TestGeminiSessionID — geminiSessionID가 RFC 4122 v4 UUID 포맷(8-4-4-4-12 하이픈, version nibble '4', variant nibble 8/9/a/b)을 내고 연속 두 호출이 서로 다른 값(랜덤성)인지 검증. 외부 의존성 없이 crypto/rand만 사용함을 행동으로 고정.

package llm

import (
	"regexp"
	"testing"
)

var uuidV4Re = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

// TestGeminiSessionID: geminiSessionID emits an RFC 4122 v4 UUID and two successive
// calls differ (randomness).
func TestGeminiSessionID(t *testing.T) {
	a, err := geminiSessionID()
	if err != nil {
		t.Fatalf("geminiSessionID error: %v", err)
	}
	if !uuidV4Re.MatchString(a) {
		t.Fatalf("uuid %q does not match v4 format", a)
	}
	b, err := geminiSessionID()
	if err != nil {
		t.Fatalf("geminiSessionID error: %v", err)
	}
	if a == b {
		t.Fatalf("two ids identical %q (want randomness)", a)
	}
}
