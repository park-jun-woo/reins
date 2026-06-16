//ff:type feature=llm type=model
//ff:what grokCLISessionFallbackCase — REINS_GROK_SESSION 폴백 테이블의 한 행(이름·설정여부·값)을 담는 테스트 케이스 구조체. 모든 행은 Stateless로 폴백됨을 기대.

package llm

// grokCLISessionFallbackCase is one row of the REINS_GROK_SESSION fallback table;
// every row is expected to fall back to Stateless.
type grokCLISessionFallbackCase struct {
	name string
	set  bool
	val  string
}
