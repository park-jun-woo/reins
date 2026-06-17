//ff:func feature=llm type=helper control=selection
//ff:what newSubprocessCLI — subprocess 백엔드 이름(claude/grok/codex/geminicli)을 해당 CLI 어댑터 생성자로 디스패치한다. FromFlag의 4개 case를 한곳으로 모아 쿼리-거부 검사(checkOptsAllowed ∅)를 공유하게 한다. 알 수 없는 이름은 nil(호출부 switch가 이미 알려진 이름만 넘긴다).

package llm

// newSubprocessCLI dispatches a subprocess backend name to its CLI adapter
// constructor. FromFlag only passes the four known names; an unknown name yields nil.
func newSubprocessCLI(backend, model string) Backend {
	switch backend {
	case "claude":
		return newClaudeCLI(model)
	case "grok":
		return newGrokCLI(model)
	case "codex":
		return newCodexCLI(model)
	case "geminicli":
		return newGeminiCLI(model)
	default:
		return nil
	}
}
