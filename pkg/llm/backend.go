//ff:type feature=llm type=model
//ff:what Backend 계약 — system+user 프롬프트로 chat completion을 수행하는 생성자(L0). 판정/래칫과 무관(권위 비대칭). 어댑터(ollama/xai/gemini)가 이를 구현한다.

package llm

// Backend performs a chat completion from a system + user prompt. It is the
// generation stage only (L0) — it has no authority over the gate's PASS lock
// (authority asymmetry).
type Backend interface {
	Complete(system, user string) (string, error)
}
