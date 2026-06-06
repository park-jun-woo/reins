//ff:type feature=llm type=adapter
//ff:what Ollama — 로컬 Ollama 서버용 Backend. Model·BaseURL(빈 값이면 localhost:11434)·NumCtx(0이면 autoNumCtx 자동 산정)를 갖는다. API 키 없음(local).

package llm

// defaultOllamaBaseURL is the local ollama server address used when BaseURL is empty.
const defaultOllamaBaseURL = "http://localhost:11434"

// Ollama is a Backend that calls a local ollama server. No API key (local).
type Ollama struct {
	Model   string
	BaseURL string // empty ⇒ "http://localhost:11434"
	NumCtx  int    // 0 ⇒ auto (autoNumCtx on system+user); else used verbatim
}
