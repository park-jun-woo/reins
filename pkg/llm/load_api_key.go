//ff:func feature=llm type=loader control=sequence level=error
//ff:what loadAPIKey — backend별 API 키를 환경변수에서만 읽는다(yaml 폴백 없음 — 선언한 의존성 경계 net/http+encoding/json을 깨지 않는다). 키 불필요 backend·미설정 env는 명확한 에러.

package llm

import (
	"fmt"
	"os"
)

// loadAPIKey returns the API key for backend from its environment variable. Unlike
// yongol, there is no XDG credentials.yaml fallback (env-only, no yaml dependency).
func loadAPIKey(backend string) (string, error) {
	envVar := backendEnvVar(backend)
	if envVar == "" {
		return "", fmt.Errorf("backend %q needs no API key", backend)
	}
	if v := os.Getenv(envVar); v != "" {
		return v, nil
	}
	return "", fmt.Errorf("load API key for %s: env %s not set", backend, envVar)
}
