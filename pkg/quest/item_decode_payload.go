//ff:func feature=quest type=helper control=sequence level=error
//ff:what Item.Payload(raw)를 v로 unmarshal해 도메인 타입으로 재수화한다. 빈 Payload는 nil(에러 없음)로 graceful.

package quest

import "encoding/json"

// DecodePayload unmarshals it.Payload into v. An empty payload is a no-op
// returning nil, so consumers get a zero-valued target rather than an error.
func (it *Item) DecodePayload(v any) error {
	if len(it.Payload) == 0 {
		return nil
	}
	return json.Unmarshal(it.Payload, v)
}
