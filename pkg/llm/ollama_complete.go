//ff:func feature=llm type=adapter control=sequence level=error
//ff:what Ollama.Complete — 로컬 Ollama 서버(POST {BaseURL}/api/chat)에 chat completion 요청. temperature 0 고정, num_predict 2048, num_ctx는 NumCtx!=0이면 그 값, 0이면 autoNumCtx(system+user) 자동 산정. BaseURL 기본 "http://localhost:11434". 전송·파싱은 doOllamaRequest에 위임.

package llm

// Complete posts a chat completion to {BaseURL}/api/chat and returns message.content.
// temperature is fixed at 0; num_ctx is auto-sized from the prompt when NumCtx == 0.
func (o Ollama) Complete(system, user string) (string, error) {
	base := o.BaseURL
	if base == "" {
		base = defaultOllamaBaseURL
	}
	numCtx := o.NumCtx
	if numCtx == 0 {
		numCtx = autoNumCtx(system + user)
	}
	body := map[string]any{
		"model": o.Model,
		"messages": []map[string]string{
			{"role": "system", "content": system},
			{"role": "user", "content": user},
		},
		"stream": false,
		"options": map[string]any{
			"temperature": 0,
			"num_predict": numPredict,
			"num_ctx":     numCtx,
		},
	}
	return doOllamaRequest(base+"/api/chat", body)
}
