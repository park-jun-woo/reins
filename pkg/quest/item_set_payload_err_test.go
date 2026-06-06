//ff:func feature=quest type=helper control=sequence level=error
//ff:what marshal 불가 값(채널)을 SetPayload에 넘기면 에러를 반환하고 Payload를 더럽히지 않는지 검증한다.

package quest

import "testing"

// TestSetPayloadMarshalError: a value json.Marshal cannot encode (a channel)
// makes SetPayload return the marshal error and leaves Payload untouched.
func TestSetPayloadMarshalError(t *testing.T) {
	it := &Item{Key: "a", State: TODO}
	if err := it.SetPayload(make(chan int)); err == nil {
		t.Fatal("SetPayload(chan): want marshal error, got nil")
	}
	if it.Payload != nil {
		t.Fatalf("SetPayload error path: want Payload nil, got %q", it.Payload)
	}
}
