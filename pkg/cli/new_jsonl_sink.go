//ff:func feature=cli type=helper control=sequence level=error
//ff:what JSONL export sink을 생성한다. path의 부모 디렉터리가 있으면 미리 만들어 둔다(append 대상 파일은 첫 Emit 때 생성).

package cli

import (
	"os"
	"path/filepath"
)

// newJSONLSink returns a sink writing to path, creating the parent directory.
func newJSONLSink(path string) (*jsonlSink, error) {
	if dir := filepath.Dir(path); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, err
		}
	}
	return &jsonlSink{path: path}, nil
}
