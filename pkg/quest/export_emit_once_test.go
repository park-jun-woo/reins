//ff:func feature=quest type=helper control=sequence
//ff:what Export가 terminal·미방출 아이템만 방출하고(비종단·기방출은 건너뜀) 재호출 시 0건인지(emit-once 래칫) 검증한다.

package quest

import "testing"

func TestExportEmitOnce(t *testing.T) {
	s := &Session{Items: []*Item{
		{Key: "todo", State: TODO},               // skipped: non-terminal
		{Key: "pass", State: PASS},               // emitted
		{Key: "review", State: REVIEW},           // emitted
		{Key: "done", State: DONE},               // emitted
		{Key: "skip", State: SKIPPED},            // emitted
		{Key: "block", State: BLOCKED},           // emitted
		{Key: "old", State: PASS, Emitted: true}, // skipped: already emitted
	}}

	sink := &memSink{}
	n, err := Export(s, sink)
	if err != nil {
		t.Fatalf("Export error: %v", err)
	}
	if n != 5 {
		t.Fatalf("first Export n = %d, want 5", n)
	}
	if len(sink.keys) != 5 {
		t.Fatalf("sink got %d keys, want 5: %v", len(sink.keys), sink.keys)
	}

	n2, err := Export(s, sink)
	if err != nil {
		t.Fatalf("second Export error: %v", err)
	}
	if n2 != 0 || len(sink.keys) != 5 {
		t.Errorf("re-export n = %d, total keys = %d, want 0 new (emit-once)", n2, len(sink.keys))
	}
}
