//ff:func feature=llm type=loader control=iteration dimension=1
//ff:what TestBackendEnvVar — 각 backend가 API 키 env var 이름으로 매핑되는지 테이블로 검증(xai→XAI_API_KEY, gemini→GEMINI_API_KEY, ollama·unknown→빈 문자열).

package llm

import (
	"testing"
)

// TestBackendEnvVar maps each backend to its API-key env var name.
func TestBackendEnvVar(t *testing.T) {
	cases := map[string]string{
		"xai":     "XAI_API_KEY",
		"gemini":  "GEMINI_API_KEY",
		"ollama":  "", // local, no key
		"unknown": "",
	}
	for backend, want := range cases {
		if got := backendEnvVar(backend); got != want {
			t.Fatalf("backendEnvVar(%q) = %q, want %q", backend, got, want)
		}
	}
}
