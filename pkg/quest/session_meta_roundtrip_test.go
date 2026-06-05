//ff:func feature=quest type=helper control=sequence level=error
//ff:what Save→Load 라운드트립에서 Session.Meta가 보존되는지 검증한다(G2 영속성).

package quest

import (
	"path/filepath"
	"testing"
)

// TestSessionMetaRoundTrip: Meta survives the Save/Load JSON round-trip.
func TestSessionMetaRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "session.json")

	want := New()
	want.SetMeta("ua", "ccnews-bot/1.0")
	want.SetMeta("cursor", "2026-06-05")
	if err := want.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if ua, ok := got.GetMeta("ua"); !ok || ua != "ccnews-bot/1.0" {
		t.Fatalf("GetMeta(ua) = %v, %v; want ccnews-bot/1.0, true", ua, ok)
	}
	if cur, ok := got.GetMeta("cursor"); !ok || cur != "2026-06-05" {
		t.Fatalf("GetMeta(cursor) = %v, %v; want 2026-06-05, true", cur, ok)
	}
}
