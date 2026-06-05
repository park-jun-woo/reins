//ff:func feature=ground type=helper control=sequence
//ff:what Snapshot.MXResolves 1회-resolve 증명 — 같은 도메인을 두 번 읽어도 bool이 같고 LookupMX는 정확히 1회만 호출됨을 fakeResolver 호출 카운터로 검증한다(키당 최대 1회 resolve).

package ground

import "testing"

// TestSnapshotMXResolvesCachesOnce proves the snapshot resolves a domain at most
// once: two reads of the same domain yield the same bool with exactly one LookupMX
// call.
func TestSnapshotMXResolvesCachesOnce(t *testing.T) {
	fr := &fakeResolver{mx: map[string]bool{"example.com": true}}
	s := NewSnapshot(fr)

	ok1, err1 := s.MXResolves("example.com")
	ok2, err2 := s.MXResolves("example.com")
	if err1 != nil || err2 != nil {
		t.Fatalf("unexpected error: %v / %v", err1, err2)
	}
	if !ok1 || !ok2 {
		t.Fatalf("want deliverable, got %v %v", ok1, ok2)
	}
	if fr.lookupCalls != 1 {
		t.Fatalf("want 1 lookup call, got %d", fr.lookupCalls)
	}
}
