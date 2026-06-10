//ff:type feature=cli type=model
//ff:func feature=cli type=helper control=sequence level=error
//ff:what 테스트용 fake sink(failAtSink). n번째 Emit(1-based)에서 에러를 주입해 exportAndSave가 부분 실패 시에도 세션을 저장하고 Emitted 래칫을 디스크에 보존하는지, 재export 시 1번째가 재방출되지 않는지(emit-once 보존 게이트)를 검증한다.

package cli

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/reins/pkg/quest"
)

// failAtSink fails on the nth Emit (1-based) and records emitted keys otherwise.
type failAtSink struct {
	calls  int
	failAt int
	keys   []string
}

func (f *failAtSink) Emit(it *quest.Item) error {
	f.calls++
	if f.calls == f.failAt {
		return errors.New("emit boom")
	}
	f.keys = append(f.keys, it.Key)
	return nil
}

// TestExportAndSavePartialFailure: with two terminal items, the first Emit succeeds
// and the second fails — exportAndSave must still save the session so the first
// item's Emitted ratchet persists, and a re-export after reload emits only the
// second item (the first is never re-emitted: emit-once preserved).
func TestExportAndSavePartialFailure(t *testing.T) {
	session := filepath.Join(t.TempDir(), "session.json")
	s := &quest.Session{Version: 1, Items: []*quest.Item{
		{Key: "a", State: quest.PASS},
		{Key: "b", State: quest.PASS},
	}}

	n, err := exportAndSave(s, &failAtSink{failAt: 2}, session)
	if err == nil {
		t.Fatal("exportAndSave = nil error, want emit error")
	}
	if n != 1 {
		t.Fatalf("n = %d, want 1 (first emitted before failure)", n)
	}

	// The partial ratchet must be on disk despite the export error.
	reloaded, err := quest.Load(session)
	if err != nil {
		t.Fatalf("reload session: %v", err)
	}
	a, err := reloaded.Find("a")
	if err != nil {
		t.Fatal(err)
	}
	if !a.Emitted {
		t.Fatal("item a Emitted ratchet lost: session was not saved after the partial export")
	}

	// Re-export must emit only "b" — "a" is never re-emitted.
	sink := &failAtSink{failAt: 0}
	n2, err := exportAndSave(reloaded, sink, session)
	if err != nil {
		t.Fatalf("re-export: %v", err)
	}
	if n2 != 1 || len(sink.keys) != 1 || sink.keys[0] != "b" {
		t.Fatalf("re-export n=%d keys=%v, want exactly [b] (no duplicate emission of a)", n2, sink.keys)
	}
}
