//ff:type feature=quest type=model
//ff:what export 수신 인터페이스. 포맷(JSONL/CSV 등)은 구현이 선택한다(cli가 JSONL sink 제공).

package quest

// Sink receives terminal item records during export. Implementations choose the
// format (JSONL, CSV, …); reins ships a JSONL file sink in package cli.
type Sink interface {
	Emit(it *Item) error
}
