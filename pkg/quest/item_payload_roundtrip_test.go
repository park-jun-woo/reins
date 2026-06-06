//ff:func feature=quest type=helper control=sequence level=error
//ff:what SetPayload→Save→Load→DecodePayload가 무손실 라운드트립인지 검증한다(G2 영속성).

package quest

import (
	"path/filepath"
	"testing"
)

// TestPayloadRoundTrip: a concrete payload set via SetPayload survives
// Save→Load and decodes back to an identical value via DecodePayload — the
// lossless rehydration that a raw json.RawMessage guarantees but `any` (which
// loads as map[string]any) does not.
func TestPayloadRoundTrip(t *testing.T) {
	want := payloadDoc{URL: "https://e.com/a", Lang: "ko"}
	it := &Item{Key: "a", State: TODO}
	if err := it.SetPayload(&want); err != nil {
		t.Fatalf("SetPayload: %v", err)
	}
	s := &Session{Version: 1, Items: []*Item{it}}

	path := filepath.Join(t.TempDir(), "session.json")
	if err := s.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}
	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	var got payloadDoc
	if err := loaded.Items[0].DecodePayload(&got); err != nil {
		t.Fatalf("DecodePayload: %v", err)
	}
	if got != want {
		t.Fatalf("round trip: got %+v, want %+v", got, want)
	}
}
