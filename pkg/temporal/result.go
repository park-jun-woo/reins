//ff:type feature=temporal type=model
//ff:what Resolve의 산출 결과. Value=정규 그레고리력 ISO(단일 날짜, 또는 IsInterval일 때 "start/end"), Determined=명세를 결정론적으로 환원했는지(false면 규칙이 REVIEW로 매핑).

package temporal

// Result is the outcome of Resolve: the normalized Gregorian ISO Value (a single
// date, or "start/end" when IsInterval), and whether the spec could be Determined.
type Result struct {
	Value      string `json:"value,omitempty"`
	IsInterval bool   `json:"is_interval,omitempty"`
	Determined bool   `json:"determined"`
}
