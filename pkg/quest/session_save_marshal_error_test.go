//ff:func feature=quest type=helper control=sequence level=error
//ff:what Save가 직렬화 불가 Payload(채널)에서 파일 쓰기 전에 marshal 에러를 내는지 검증한다.

package quest

import (
	"path/filepath"
	"testing"
)

func TestSaveMarshalError(t *testing.T) {
	// An unmarshalable Payload (a channel) makes json.MarshalIndent fail before any
	// file write is attempted.
	s := &Session{Version: 1, Items: []*Item{{Key: "a", State: TODO, Payload: make(chan int)}}}
	path := filepath.Join(t.TempDir(), "session.json")
	if err := s.Save(path); err == nil {
		t.Fatal("Save unmarshalable payload: want error, got nil")
	}
}
