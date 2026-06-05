//ff:type feature=gate type=model
//ff:what 규칙의 자기문서화 카탈로그 항목(자동 rulebook). ID·Level·Desc — cli `rules`가 출력하는 "이 게이트가 막는 치즈 목록"의 한 행.

package gate

// RuleMeta is the self-documenting catalog entry for a rule (the auto rulebook).
type RuleMeta struct {
	ID    string `json:"id"`
	Level Level  `json:"level"`
	Desc  string `json:"desc"`
}
