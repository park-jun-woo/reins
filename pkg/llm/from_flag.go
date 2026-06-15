//ff:func feature=llm type=helper control=selection
//ff:what FromFlag — "backend:model" 문자열을 Backend로 만든다(yongol parseModelFlag 이식). 첫 ':' 앞이 backend, 나머지 전부가 model(ollama:gemma4:e4b → Ollama{Model:"gemma4:e4b"}). ollama BaseURL은 env REINS_OLLAMA_URL 오버라이드. xai→OpenAICompat(x.ai endpoint), gemini→Gemini. 지원 외 backend·빈 model은 에러.

package llm

import (
	"fmt"
	"os"
	"strings"
)

// FromFlag turns a "backend:model" string into a Backend. The first ':' separates
// the backend from the model, so the model may itself contain colons (ollama
// "gemma4:e4b"). ollama reads REINS_OLLAMA_URL to override its BaseURL.
func FromFlag(flag string) (Backend, error) {
	idx := strings.Index(flag, ":")
	if idx < 0 {
		return nil, fmt.Errorf("invalid --model %q: expected format backend:model (e.g. ollama:gemma4:e4b)", flag)
	}
	backend := flag[:idx]
	model := flag[idx+1:]
	if model == "" {
		return nil, fmt.Errorf("invalid --model %q: model name is empty", flag)
	}
	switch backend {
	case "ollama":
		return Ollama{Model: model, BaseURL: os.Getenv("REINS_OLLAMA_URL")}, nil
	case "xai":
		return OpenAICompat{URL: "https://api.x.ai/v1/chat/completions", Backend: backend, Model: model}, nil
	case "gemini":
		return Gemini{Model: model}, nil
	case "claude":
		return newClaudeCLI(model), nil
	default:
		return nil, fmt.Errorf("invalid --model backend %q: supported backends: ollama, xai, gemini, claude", backend)
	}
}
