//ff:func feature=ground type=helper control=sequence
//ff:what Snapshot.HTTPBody 에러 캐시 증명 — 실패한 fetch의 에러도 캐시되어, 두 번째 읽기에서 resolver 추가 호출 없이 같은 에러가 재현됨을 fakeResolver 호출 카운터로 검증한다.

package ground

import (
	"errors"
	"testing"
)

// TestSnapshotHTTPBodyCachesError proves a failed fetch is cached: the error is
// reproduced on a second read without a second resolver call.
func TestSnapshotHTTPBodyCachesError(t *testing.T) {
	want := errors.New("boom")
	fr := &fakeResolver{fetchErr: map[string]error{"http://x": want}}
	s := NewSnapshot(fr)

	_, err1 := s.HTTPBody("http://x")
	_, err2 := s.HTTPBody("http://x")
	if !errors.Is(err1, want) || !errors.Is(err2, want) {
		t.Fatalf("want cached error, got %v / %v", err1, err2)
	}
	if fr.fetchCalls != 1 {
		t.Fatalf("want 1 fetch call, got %d", fr.fetchCalls)
	}
}
