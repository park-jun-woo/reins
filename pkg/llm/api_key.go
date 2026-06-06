//ff:func feature=llm type=loader control=selection
//ff:what backendEnvVar — backend별 API 키를 담은 환경변수 이름을 반환한다(yongol backendEnvVar 규칙: xai→XAI_API_KEY, gemini→GEMINI_API_KEY). ollama는 로컬이라 키 불필요(빈 문자열).

package llm

// backendEnvVar returns the environment variable name holding the API key for a
// backend (yongol backendEnvVar rule). ollama needs none (local).
func backendEnvVar(backend string) string {
	switch backend {
	case "xai":
		return "XAI_API_KEY"
	case "gemini":
		return "GEMINI_API_KEY"
	default:
		return ""
	}
}
