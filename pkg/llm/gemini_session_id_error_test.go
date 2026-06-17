//ff:func feature=llm type=helper control=sequence level=error
//ff:what TestGeminiSessionIDError — geminiRandRead seam이 실패하면 geminiSessionID가 id 없이 에러를 표면화하는지 검증(실 crypto/rand는 사실상 실패 안 하므로 seam 주입으로 에러 경로를 고정).

package llm

import (
	"errors"
	"testing"
)

// TestGeminiSessionIDError: a failing randomness seam surfaces as an error (no id).
func TestGeminiSessionIDError(t *testing.T) {
	orig := geminiRandRead
	defer func() { geminiRandRead = orig }()
	geminiRandRead = func(b []byte) (int, error) { return 0, errors.New("no entropy") }

	if _, err := geminiSessionID(); err == nil {
		t.Fatal("geminiSessionID = nil error, want error from failing reader")
	}
}
