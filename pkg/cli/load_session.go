//ff:func feature=cli type=helper control=sequence level=error
//ff:what sessionPath의 세션을 로드한다. 파일이 없으면 빈 세션(quest.New())을 새로 만들어 반환한다 — 첫 scan을 위해 부재를 에러로 보지 않는다.

package cli

import (
	"os"

	"github.com/park-jun-woo/reins/pkg/quest"
)

// loadSession returns the session at path, creating a fresh one if absent.
func loadSession(path string) (*quest.Session, error) {
	s, err := quest.Load(path)
	if os.IsNotExist(err) {
		return quest.New(), nil
	}
	return s, err
}
