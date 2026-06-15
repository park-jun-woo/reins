//ff:func feature=llm type=helper control=sequence
//ff:what withNoToolsPreamble — claudeNoToolsPreamble(도구없음·주어진 텍스트로만·페이로드만) 상수를 소비자 system 앞에 선행 결합한다. 빈 system이면 프리앰블 단독 반환. --tools ""(런타임 도구 차단)와 짝을 이뤄 도구-유발 프롬프트에서 모델이 도구호출을 텍스트로 뱉지 않게 한다(BUG-001).

package llm

// claudeNoToolsPreamble pins the L0-generator contract into claude's system prompt:
// no tools, work only from the given text, emit only the payload. It pairs with
// --tools "" (runtime tool disable) so the model neither attempts a tool nor narrates
// a tool call as text when a prompt invites file access (BUG-001).
const claudeNoToolsPreamble = "You have NO tools. You cannot read files, run commands, " +
	"browse, or access anything beyond this prompt. Work ONLY from the text given here. " +
	"Output ONLY the requested payload in the exact format the prompt specifies — " +
	"no prose, no tool calls, no commands, no explanation. If required information is " +
	"absent from the given text, emit the payload that the format prescribes for that case."

// withNoToolsPreamble prepends the fixed no-tools preamble to the consumer system
// prompt. An empty system yields the preamble alone.
func withNoToolsPreamble(system string) string {
	if system == "" {
		return claudeNoToolsPreamble
	}
	return claudeNoToolsPreamble + "\n\n" + system
}
