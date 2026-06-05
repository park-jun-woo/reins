//ff:func feature=quest type=helper control=sequence level=error
//ff:what Export가 sink 에러 시 그때까지 방출한 건수와 함께 중단하고, 실패 아이템은 Emitted로 표시하지 않는지 검증한다.

package quest

import "testing"

// TestExportSinkError: a sink error aborts export, returning the count emitted so far
// and not marking the failed item as Emitted.
func TestExportSinkError(t *testing.T) {
	first := &Item{Key: "a", State: PASS}
	second := &Item{Key: "b", State: PASS}
	s := &Session{Items: []*Item{first, second}}

	n, err := Export(s, &failSink{failAt: 2})
	if err == nil {
		t.Fatal("Export: want sink error, got nil")
	}
	if n != 1 {
		t.Errorf("n = %d, want 1 (first emitted before failure)", n)
	}
	if !first.Emitted {
		t.Error("first item should be marked Emitted")
	}
	if second.Emitted {
		t.Error("failed item must not be marked Emitted")
	}
}
