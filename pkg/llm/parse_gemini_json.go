//ff:func feature=llm type=adapter control=sequence level=error
//ff:what parseGeminiJSON — `gemini -o json` 봉투 {response, stats{tools{calls}}, error}에서 response를 뽑는다. error 필드가 non-null이면 그 봉투를 에러로 표면화한다. json 봉투가 아니면(알려진 이슈 #11184: 일부 버전이 -o json을 미준수하고 plain 텍스트를 뱉음) 폴백으로 trimmed plain stdout을 response로 취한다(완전히 비어 있으면 에러). stats.tools.calls는 스모크 단계(인증 게이트)에서 단일샷 검증용이라 여기선 추출하지 않는다(결과=response).

package llm

import (
	"encoding/json"
	"fmt"
	"strings"
)

// parseGeminiJSON extracts the `response` field from the `gemini -o json` envelope
// {response, stats, error}. A non-null `error` field surfaces as a Complete error.
// If stdout is not the json envelope at all (known issue #11184: some builds ignore
// -o json and emit plain text), it falls back to the trimmed plain stdout as the
// response; a completely empty output is an error.
func parseGeminiJSON(stdout string) (string, error) {
	var env struct {
		Response string          `json:"response"`
		Error    json.RawMessage `json:"error"`
	}
	if err := json.Unmarshal([]byte(stdout), &env); err != nil {
		// Fallback (#11184): no json envelope — take the plain text as the response.
		text := strings.TrimSpace(stdout)
		if text == "" {
			return "", fmt.Errorf("gemini: empty output, no json envelope")
		}
		return text, nil
	}
	if len(env.Error) > 0 && string(env.Error) != "null" {
		return "", fmt.Errorf("gemini error envelope: %s", strings.TrimSpace(string(env.Error)))
	}
	return env.Response, nil
}
