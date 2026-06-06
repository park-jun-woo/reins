//ff:func feature=quest type=helper control=sequence level=error
//ff:what 빈 Payload에 대한 DecodePayload가 no-op(에러 없음·zero-value 유지)인지 검증한다.

package quest

import "testing"

// TestPayloadEmptyNilSafe: DecodePayload on an item with no payload is a no-op
// that leaves the target zero-valued and returns nil (no error).
func TestPayloadEmptyNilSafe(t *testing.T) {
	it := &Item{Key: "a", State: TODO}
	var got payloadDoc
	if err := it.DecodePayload(&got); err != nil {
		t.Fatalf("DecodePayload empty: want nil error, got %v", err)
	}
	if (got != payloadDoc{}) {
		t.Fatalf("DecodePayload empty: want zero value, got %+v", got)
	}
}
