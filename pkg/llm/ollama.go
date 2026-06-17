//ff:type feature=llm type=adapter
//ff:what Ollama — 로컬 Ollama 서버용 Backend. Model·BaseURL(빈 값이면 localhost:11434)·NumCtx(0이면 autoNumCtx 자동 산정)·MaxOutputTokens(0⇒2048, options.num_predict)·Temperature(nil⇒0)·Think(nil⇒모델기본, false⇒options.think=false)를 갖는다. API 키 없음(local).

package llm

// defaultOllamaBaseURL is the local ollama server address used when BaseURL is empty.
const defaultOllamaBaseURL = "http://localhost:11434"

// Ollama is a Backend that calls a local ollama server. No API key (local).
type Ollama struct {
	Model           string
	BaseURL         string   // empty ⇒ "http://localhost:11434"
	NumCtx          int      // 0 ⇒ auto (autoNumCtx on system+user); else used verbatim
	MaxOutputTokens int      // 0 ⇒ 2048; maps to options.num_predict
	Temperature     *float64 // nil ⇒ 0 (current); else used verbatim
	Think           *bool    // nil ⇒ model default; false ⇒ options.think=false
}
