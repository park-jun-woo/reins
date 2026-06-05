//ff:func feature=quest type=helper control=sequence level=error
//ff:what meta 필드 없는 구버전 세션 JSON도 역직렬화되고 GetMeta가 nil-safe한지 검증한다(후방호환).

package quest

import (
	"os"
	"path/filepath"
	"testing"
)

// TestSessionMetaLegacyLoad: a pre-Meta session JSON loads cleanly and GetMeta
// is nil-safe (the field is simply absent).
func TestSessionMetaLegacyLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "session.json")
	legacy := `{"version":1,"items":[{"key":"a","state":"todo"}]}`
	if err := os.WriteFile(path, []byte(legacy), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	got, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.Meta != nil {
		t.Fatalf("Meta = %v, want nil for legacy session", got.Meta)
	}
	if v, ok := got.GetMeta("missing"); ok || v != nil {
		t.Fatalf("GetMeta on nil Meta = %v, %v; want nil, false", v, ok)
	}
}
