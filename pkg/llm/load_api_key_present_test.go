//ff:func feature=llm type=loader control=sequence level=error
//ff:what TestLoadAPIKeyPresent — backend의 env var가 설정돼 있으면 그 키를 읽어 반환하는지 검증.

package llm

import (
	"testing"
)

// TestLoadAPIKeyPresent: the key is read from the backend's env var.
func TestLoadAPIKeyPresent(t *testing.T) {
	t.Setenv("XAI_API_KEY", "secret-xai")
	got, err := loadAPIKey("xai")
	if err != nil {
		t.Fatalf("loadAPIKey error: %v", err)
	}
	if got != "secret-xai" {
		t.Fatalf("key = %q, want secret-xai", got)
	}
}
