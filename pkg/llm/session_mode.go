//ff:type feature=llm type=model
//ff:what SessionMode — sessionKind를 감싸 ClaudeCLI가 bare enum 대신 명명 구조체로 세션 동작을 설정하게 한다. Kind 필드 하나(Stateless/Continue).

package llm

// SessionMode wraps a sessionKind so callers configure the session behaviour
// through a named struct rather than a bare enum.
type SessionMode struct {
	Kind sessionKind
}
