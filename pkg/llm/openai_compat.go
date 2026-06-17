//ff:type feature=llm type=adapter
//ff:what OpenAICompat — OpenAI 호환 chat completions endpoint(xai 등)용 Backend. URL·Backend(키 조회용 env 이름)·Model·MaxOutputTokens(0⇒2048, max_tokens)·Temperature(nil⇒0)을 갖는다. Backend가 API 키 env var 이름을 가리킨다(예: "xai" ⇒ XAI_API_KEY).

package llm

// OpenAICompat is a Backend for an OpenAI-compatible chat completions endpoint (xai).
// Backend names the env var holding the API key (e.g. "xai" ⇒ XAI_API_KEY).
type OpenAICompat struct {
	URL             string
	Backend         string
	Model           string
	MaxOutputTokens int      // 0 ⇒ 2048; maps to max_tokens
	Temperature     *float64 // nil ⇒ 0 (current); else used verbatim
}
