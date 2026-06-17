//ff:type feature=llm type=model
//ff:what backendOpts — `--model` 쿼리(`?k=v&…`)에서 파싱한 타입드 백엔드 옵션. MaxOutputTokens·NumCtx(int), Temperature·Think(포인터로 미설정 구분), present(실제로 주어진 canonical 키 집합 — int zero값이 "미설정 vs 0"을 못 가리므로 FromFlag가 백엔드별 허용 키와 대조해 미적용 키를 loud 거부하는 데 쓴다).

package llm

// backendOpts holds typed options parsed from a `--model` query string. present
// records which canonical keys were actually supplied so FromFlag can reject keys
// the chosen backend does not support (an int zero value cannot distinguish
// "unset" from "explicitly 0").
type backendOpts struct {
	MaxOutputTokens int
	NumCtx          int
	Temperature     *float64
	Think           *bool
	present         map[string]bool
}
