//ff:func feature=llm type=helper control=sequence
//ff:what splitModelQuery — FromFlag가 backend|model로 가른 뒤의 model 부분을 첫 '?' 기준으로 model|raw(쿼리)로 재분리한다. '?' 없으면 raw는 빈 문자열(현행 동치). model은 콜론을 포함할 수 있다(ollama gemma4:e4b) — '?' 앞 전체가 model.

package llm

import "strings"

// splitModelQuery splits the model part (everything after the backend's first ':')
// on the first '?' into the model name and the raw query string. With no '?' the
// raw string is empty, preserving the pre-Phase017 behavior byte-for-byte.
func splitModelQuery(modelPart string) (model string, raw string) {
	idx := strings.Index(modelPart, "?")
	if idx < 0 {
		return modelPart, ""
	}
	return modelPart[:idx], modelPart[idx+1:]
}
