//ff:type feature=llm type=model
//ff:what claudeCLISessionCase — REINS_CLAUDE_SESSION 매핑 테이블의 한 행(이름·설정여부·값·기대 sessionKind)을 담는 테스트 케이스 구조체.

package llm

// claudeCLISessionCase is one row of the REINS_CLAUDE_SESSION mapping table.
type claudeCLISessionCase struct {
	name string
	set  bool
	val  string
	want sessionKind
}
