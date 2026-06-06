//ff:type feature=llm type=model
//ff:what CallFunc — ground.Resolver와 동형의 테스트 주입용 함수형 Backend. 래핑한 함수에 위임해 Backend를 만족한다.

package llm

// CallFunc is a function-typed Backend for test injection (homologous to
// ground.Resolver injection).
type CallFunc func(system, user string) (string, error)
