//ff:func feature=quest type=helper control=sequence level=error
//ff:what 도메인 값을 JSON으로 marshal해 Item.Payload(raw)에 스냅샷으로 담는다. 변이를 모두 끝낸 뒤 호출해야 변이가 보존된다.

package quest

import "encoding/json"

// SetPayload marshals v to JSON and stores the raw bytes in it.Payload. It is a
// snapshot: callers must finish all mutations of v before calling, otherwise
// later changes are not reflected in the persisted payload.
func (it *Item) SetPayload(v any) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	it.Payload = b
	return nil
}
