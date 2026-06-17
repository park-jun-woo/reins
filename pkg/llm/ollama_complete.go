//ff:func feature=llm type=adapter control=sequence level=error
//ff:what Ollama.Complete — 로컬 Ollama 서버(POST {BaseURL}/api/chat)에 chat completion 요청. 출력 한도 eff=MaxOutputTokens(0⇒numPredict 2048)를 options.num_predict와 autoNumCtx 예비분 양쪽에 넘겨 출력 한도 상향 시 num_ctx도 따라 오르게 한다(Phase017 A). temperature는 Temperature(nil⇒0), Think!=nil이면 options.think 반영. num_ctx는 NumCtx!=0이면 그 값, 0이면 autoNumCtx(system+user, eff). BaseURL 기본 "http://localhost:11434". 전송·파싱은 doOllamaRequest에 위임.

package llm

// Complete posts a chat completion to {BaseURL}/api/chat and returns message.content.
// eff is the effective output budget (MaxOutputTokens, 0 ⇒ numPredict 2048); it feeds
// both options.num_predict and the autoNumCtx reserve so a raised output limit also
// raises num_ctx (no silent truncation). temperature defaults to 0 (Temperature nil),
// and options.think is set only when Think is non-nil.
func (o Ollama) Complete(system, user string) (string, error) {
	base := o.BaseURL
	if base == "" {
		base = defaultOllamaBaseURL
	}
	eff := o.MaxOutputTokens
	if eff == 0 {
		eff = numPredict
	}
	numCtx := o.NumCtx
	if numCtx == 0 {
		numCtx = autoNumCtx(system+user, eff)
	}
	temperature := 0.0
	if o.Temperature != nil {
		temperature = *o.Temperature
	}
	options := map[string]any{
		"temperature": temperature,
		"num_predict": eff,
		"num_ctx":     numCtx,
	}
	if o.Think != nil {
		options["think"] = *o.Think
	}
	body := map[string]any{
		"model": o.Model,
		"messages": []map[string]string{
			{"role": "system", "content": system},
			{"role": "user", "content": user},
		},
		"stream":  false,
		"options": options,
	}
	return doOllamaRequest(base+"/api/chat", body)
}
