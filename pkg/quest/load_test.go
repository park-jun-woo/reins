//ff:func feature=quest type=helper control=sequence level=error
//ff:what Load가 Save한 유효 JSON을 Session으로 라운드트립 역직렬화하는지 검증한다.

package quest

import (
	"path/filepath"
	"testing"
)

func TestLoadRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "session.json")
	want := &Session{Version: 1, Items: []*Item{{Key: "a", State: TODO}}}
	if err := want.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.Version != 1 || len(got.Items) != 1 || got.Items[0].Key != "a" {
		t.Fatalf("Load = %+v, want round-trip of %+v", got, want)
	}
}
