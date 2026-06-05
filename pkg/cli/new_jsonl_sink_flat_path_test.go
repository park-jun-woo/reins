//ff:func feature=cli type=helper control=sequence level=error
//ff:what newJSONLSink이 단순 파일명(dir ".")엔 디렉터리 생성 없이 성공하는지 검증한다.

package cli

import "testing"

// TestNewJSONLSinkFlatPath: a bare filename (dir ".") needs no directory creation.
func TestNewJSONLSinkFlatPath(t *testing.T) {
	if _, err := newJSONLSink("out.jsonl"); err != nil {
		t.Fatalf("newJSONLSink flat: %v", err)
	}
}
