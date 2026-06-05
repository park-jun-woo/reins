//ff:func feature=cli type=helper control=sequence level=error
//ff:what Emit이 직렬화 불가 Payload(채널)에서 에러를 내는지 검증한다.

package cli

import (
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/reins/pkg/quest"
)

// TestEmitMarshalError: Emit returns an error when the item cannot be JSON-encoded
// (an unmarshalable Payload such as a channel).
func TestEmitMarshalError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "out.jsonl")
	sink, err := newJSONLSink(path)
	if err != nil {
		t.Fatalf("newJSONLSink: %v", err)
	}
	it := &quest.Item{Key: "a", State: quest.PASS, Payload: make(chan int)}
	if err := sink.Emit(it); err == nil {
		t.Fatal("Emit unmarshalable payload: want error, got nil")
	}
}
