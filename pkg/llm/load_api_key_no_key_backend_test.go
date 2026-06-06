//ff:func feature=llm type=loader control=sequence level=error
//ff:what TestLoadAPIKeyNoKeyBackend — 키가 불필요한 backend(ollama)는 에러를 반환하는지 검증(env-only 설계).

package llm

import (
	"testing"
)

// TestLoadAPIKeyNoKeyBackend: a backend that needs no key errors (env-only design).
func TestLoadAPIKeyNoKeyBackend(t *testing.T) {
	if _, err := loadAPIKey("ollama"); err == nil {
		t.Fatal("loadAPIKey(ollama) = nil error, want 'needs no API key'")
	}
}
