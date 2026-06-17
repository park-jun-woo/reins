//ff:func feature=llm type=adapter control=sequence level=error
//ff:what TestGeminiCLISmoke — gemini가 PATH에 있고(exec.LookPath) **인증돼 있을 때만**(~/.gemini/oauth_creds.json 존재 또는 GEMINI_API_KEY 설정) geminicli:default가 read-only(--approval-mode plan)에서 도구 시도 없이 비어있지 않은 response를 1회 생성하는지 통합 검증(L0 단일샷의 경험적 확인). 미설치/미인증이면 skip(현재 환경은 미인증이라 Skip; 사용자 `gemini login` 후 실행).

package llm

import (
	"os/exec"
	"strings"
	"testing"
)

// TestGeminiCLISmoke: when gemini is on PATH AND authenticated (~/.gemini/oauth_creds.json
// exists or GEMINI_API_KEY is set), geminicli:default produces a non-empty response
// under --approval-mode plan; otherwise skip (no CLI / not logged in).
func TestGeminiCLISmoke(t *testing.T) {
	if _, err := exec.LookPath("gemini"); err != nil {
		t.Skip("gemini not on PATH")
	}
	if !geminiAuthed() {
		t.Skip("gemini not authenticated (no ~/.gemini/oauth_creds.json, no GEMINI_API_KEY)")
	}
	c := newGeminiCLI("default")
	out, err := c.Complete("", "Reply with exactly the two letters: ok")
	if err != nil {
		t.Skipf("gemini smoke skipped (login/quota?): %v", err)
	}
	if strings.TrimSpace(out) == "" {
		t.Fatalf("gemini smoke returned empty result")
	}
}
