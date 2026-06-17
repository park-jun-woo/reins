//ff:func feature=llm type=helper control=selection level=error
//ff:what FromFlag — "backend:model[?k=v&…]" 문자열을 Backend로 만든다. 첫 ':'로 backend|model을 가른 뒤(model은 ':' 포함 가능 — ollama gemma4:e4b), model 부분을 splitModelQuery로 model|쿼리 재분리하고 parseBackendOpts로 타입드 옵션을 얻는다. 백엔드별 허용 키 집합과 present를 checkOptsAllowed로 대조해 미적용/미지 키는 거부(HTTP: ollama=max_output_tokens/num_ctx/temperature/think, xai·gemini=max_output_tokens/temperature; subprocess claude/grok/codex/geminicli=∅ ⇒ 쿼리 있으면 에러). ollama BaseURL은 env REINS_OLLAMA_URL 오버라이드. 지원 외 backend·빈 model은 에러. '?' 없으면 현행 동치.

package llm

import (
	"fmt"
	"os"
	"strings"
)

// FromFlag turns a "backend:model[?query]" string into a Backend. The first ':'
// separates the backend from the model (the model may contain colons, e.g. ollama
// "gemma4:e4b"); the model part is then split on its first '?' into the model name
// and a typed option query. Each backend validates the supplied option keys against
// its allowed set (no silent caps). With no '?' the result is identical to before.
func FromFlag(flag string) (Backend, error) {
	idx := strings.Index(flag, ":")
	if idx < 0 {
		return nil, fmt.Errorf("invalid --model %q: expected format backend:model (e.g. ollama:gemma4:e4b)", flag)
	}
	backend := flag[:idx]
	model, raw := splitModelQuery(flag[idx+1:])
	if model == "" {
		return nil, fmt.Errorf("invalid --model %q: model name is empty", flag)
	}
	opts, err := parseBackendOpts(raw)
	if err != nil {
		return nil, fmt.Errorf("invalid --model %q: %w", flag, err)
	}
	switch backend {
	case "ollama":
		if err := checkOptsAllowed(opts.present, "max_output_tokens", "num_ctx", "temperature", "think"); err != nil {
			return nil, fmt.Errorf("invalid --model %q: %w", flag, err)
		}
		return Ollama{
			Model:           model,
			BaseURL:         os.Getenv("REINS_OLLAMA_URL"),
			NumCtx:          opts.NumCtx,
			MaxOutputTokens: opts.MaxOutputTokens,
			Temperature:     opts.Temperature,
			Think:           opts.Think,
		}, nil
	case "xai":
		if err := checkOptsAllowed(opts.present, "max_output_tokens", "temperature"); err != nil {
			return nil, fmt.Errorf("invalid --model %q: %w", flag, err)
		}
		return OpenAICompat{
			URL:             "https://api.x.ai/v1/chat/completions",
			Backend:         backend,
			Model:           model,
			MaxOutputTokens: opts.MaxOutputTokens,
			Temperature:     opts.Temperature,
		}, nil
	case "gemini":
		if err := checkOptsAllowed(opts.present, "max_output_tokens", "temperature"); err != nil {
			return nil, fmt.Errorf("invalid --model %q: %w", flag, err)
		}
		return Gemini{Model: model, MaxOutputTokens: opts.MaxOutputTokens, Temperature: opts.Temperature}, nil
	case "claude", "grok", "codex", "geminicli":
		if err := checkOptsAllowed(opts.present); err != nil {
			return nil, fmt.Errorf("invalid --model %q: %w (subprocess backends take no query options)", flag, err)
		}
		return newSubprocessCLI(backend, model), nil
	default:
		return nil, fmt.Errorf("invalid --model backend %q: supported backends: ollama, xai, gemini, claude, grok, codex, geminicli", backend)
	}
}
