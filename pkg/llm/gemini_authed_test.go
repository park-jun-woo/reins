//ff:func feature=llm type=helper control=sequence level=error
//ff:what geminiAuthed — gemini CLI가 로그인된 것처럼 보이는지 반환하는 스모크 게이트 헬퍼: 프로세스에 GEMINI_API_KEY가 있거나 ~/.gemini/oauth_creds.json(Google OAuth 자격)이 존재하면 true. TestGeminiCLISmoke가 미인증 환경에서 t.Skip하도록 쓴다(인증 위임 — reins는 키를 보지 않음).

package llm

import (
	"os"
	"path/filepath"
)

// geminiAuthed reports whether the gemini CLI looks logged in: a Google OAuth creds
// file or a GEMINI_API_KEY env in the process. Used by the smoke test to skip on an
// unauthenticated environment.
func geminiAuthed() bool {
	if os.Getenv("GEMINI_API_KEY") != "" {
		return true
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return false
	}
	_, err = os.Stat(filepath.Join(home, ".gemini", "oauth_creds.json"))
	return err == nil
}
