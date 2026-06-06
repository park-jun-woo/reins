//ff:func feature=quest type=helper control=sequence level=error
//ff:what 깨진 JSON이 담긴 Payload를 DecodePayload하면 Unmarshal 에러를 그대로 surface하는지 검증한다.

package quest

import "testing"

// TestDecodePayloadUnmarshalError: a non-empty payload holding invalid JSON
// makes DecodePayload surface the json.Unmarshal error rather than swallow it.
func TestDecodePayloadUnmarshalError(t *testing.T) {
	it := &Item{Key: "a", State: TODO, Payload: []byte("{not json")}
	var got payloadDoc
	if err := it.DecodePayload(&got); err == nil {
		t.Fatal("DecodePayload(broken): want unmarshal error, got nil")
	}
}
