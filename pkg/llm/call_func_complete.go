//ff:func feature=llm type=model control=sequence
//ff:what CallFunc.Complete — 래핑한 함수에 위임해 Backend 계약을 만족한다(테스트 주입 seam).

package llm

// Complete satisfies Backend by delegating to the wrapped function.
func (f CallFunc) Complete(system, user string) (string, error) { return f(system, user) }
