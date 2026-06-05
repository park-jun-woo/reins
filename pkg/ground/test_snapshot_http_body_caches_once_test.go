//ff:func feature=ground type=helper control=sequence
//ff:what Snapshot.HTTPBody 1회-resolve 증명 — 같은 URL을 두 번 읽어도 본문이 같고 Fetch는 정확히 1회만 호출됨을 fakeResolver 호출 카운터로 검증한다(키당 최대 1회 resolve, 요청당 결정성).

package ground

import "testing"

// TestSnapshotHTTPBodyCachesOnce proves the snapshot resolves a URL at most once:
// two reads of the same URL yield the same body with exactly one Fetch call.
func TestSnapshotHTTPBodyCachesOnce(t *testing.T) {
	fr := &fakeResolver{bodies: map[string]string{"http://x": "BODY"}}
	s := NewSnapshot(fr)

	b1, err1 := s.HTTPBody("http://x")
	b2, err2 := s.HTTPBody("http://x")
	if err1 != nil || err2 != nil {
		t.Fatalf("unexpected error: %v / %v", err1, err2)
	}
	if b1 != "BODY" || b2 != "BODY" {
		t.Fatalf("body mismatch: %q %q", b1, b2)
	}
	if fr.fetchCalls != 1 {
		t.Fatalf("want 1 fetch call, got %d", fr.fetchCalls)
	}
}
