//ff:func feature=quest type=helper control=sequence level=error
//ff:what Load가 부재 파일을 os.IsNotExist로 보고하는지 검증한다.

package quest

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadMissingIsNotExist(t *testing.T) {
	_, err := Load(filepath.Join(t.TempDir(), "absent.json"))
	if !os.IsNotExist(err) {
		t.Fatalf("Load missing: err = %v, want os.IsNotExist", err)
	}
}
